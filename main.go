package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"

	pbAuth "auth-service"
	pb "function-service"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	kafkaavro "github.com/mycujoo/go-kafka-avro"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gopkg.in/couchbase/gocb.v1"
)

const (
	couchbaseServer         = "pptestserver.australiaeast.cloudapp.azure.com"
	couchbaseUsername       = "dev"
	couchbasePassword       = "Powerpilot!"
	bucketNameDefault       = "staging"
	bucketNameConnections   = "staging_connections"
	bucketNameConfig        = "staging_config"
	kafkaServer             = "pptestserver.australiaeast.cloudapp.azure.com:9092"
	schemaRegistryServer    = "http://pptestserver.australiaeast.cloudapp.azure.com:8081"
	topicFunctionDownlink   = "FUNCTION_DOWNLINK"
	serverPort              = 9070
	grpcAuthServerAddress   = "127.0.0.1:9030"
	rolePowerpilotAdmin     = "powerpilot-admin"
	rolePowerpilotInstaller = "powerpilot-installer"
)

var (
	bucketConnections *gocb.Bucket
	bucketDefault     *gocb.Bucket
	bucketConfig      *gocb.Bucket
	kafkaProducer     *kafka.Producer
	avroProducer      *kafkaavro.Producer
	grpcAuthClient    pbAuth.AuthServiceClient
)

type FunctionDetails struct {
	Index      int32  `json:"index"`
	Name       string `json:"name"`
	Param1Type string `json:"param1type"`
	Param2Type string `json:"param2type"`
	Param3Type string `json:"param3type"`
	Param4Type string `json:"param4type"`
}

type functionServiceServer struct {
	pb.UnimplementedDeviceFunctionServiceServer
}

func init() {
	cluster, err := gocb.Connect("couchbase://" + couchbaseServer)
	if err != nil {
		panic(err)
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: couchbaseUsername,
		Password: couchbasePassword,
	})
	bucketConnections, err = cluster.OpenBucket(bucketNameConnections, "")
	if err != nil {
		panic(err)
	}

	bucketDefault, err = cluster.OpenBucket(bucketNameDefault, "")
	if err != nil {
		panic(err)
	}

	bucketConfig, err = cluster.OpenBucket(bucketNameConfig, "")
	if err != nil {
		panic(err)
	}

	srClient, err := kafkaavro.NewCachedSchemaRegistryClient(schemaRegistryServer)
	if err != nil {
		panic(err)
	}

	kafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaServer,
	})
	if err != nil {
		panic(err)
	}

	contents, err := ioutil.ReadFile("schema.json")
	if err != nil {
		panic(err)
	}
	schema := string(contents)
	avroProducer, err = kafkaavro.NewProducer(
		topicFunctionDownlink,
		"\"string\"",
		schema,
		kafkaProducer,
		srClient,
	)
	if err != nil {
		panic(err)
	}
}

func getFirmwareVersion(deviceEUI string) (int32, error) {
	queryString := "SELECT b.Desired.Config.firmware FROM " + bucketNameConfig + " b " +
		"WHERE meta(b).id = $1"
	query := gocb.NewN1qlQuery(queryString)
	rows, err := bucketConfig.ExecuteN1qlQuery(query, []interface{}{
		deviceEUI,
	})
	if err != nil {
		return -1, err
	}

	var row map[string]interface{}
	if !rows.Next(&row) {
		return -1, errors.New("Could not find device " + deviceEUI)
	}

	firmware := row["firmware"].(float64)
	return int32(firmware), nil
}

func getFunctionsForFirmware(firmware int32) (*pb.DeviceFunctions, error) {
	queryString := "SELECT f.i AS `index`, f.n AS `name`, f.a AS `param1type`, " +
		"f.b AS `param2type`, f.c AS `param3type`, f.d AS `param4type` " +
		"FROM " + bucketNameDefault +
		" s UNNEST ppschema f WHERE s.Type='functionschema' " +
		"AND s.ppver = $1"

	query := gocb.NewN1qlQuery(queryString)
	rows, err := bucketDefault.ExecuteN1qlQuery(query, []interface{}{firmware})
	if err != nil {
		return nil, err
	}

	var functions []*pb.DeviceFunction
	var row map[string]interface{}
	for rows.Next(&row) {
		function := &pb.DeviceFunction{
			FunctionName: fmt.Sprintf("%v", row["name"]),
			Index:        int32(row["index"].(float64)),
			Param1Type:   fmt.Sprintf("%v", row["param1type"]),
			Param2Type:   fmt.Sprintf("%v", row["param2type"]),
			Param3Type:   fmt.Sprintf("%v", row["param3type"]),
			Param4Type:   fmt.Sprintf("%v", row["param4type"]),
		}
		functions = append(functions, function)
	}
	return &pb.DeviceFunctions{
		Functions: functions,
	}, nil
}

func getFunctionDetailsByName(name string, firmware int32) (map[string]interface{}, error) {
	queryString := "SELECT f.i AS `index`, f.n AS `name`, f.a AS `param1type`, " +
		"f.b AS `param2type`, f.c AS `param3type`, f.d AS `param4type` " +
		"FROM " + bucketNameDefault +
		" s UNNEST ppschema f WHERE s.Type='functionschema' " +
		"AND f.n = $1 " +
		"AND s.ppver = $2"
	query := gocb.NewN1qlQuery(queryString)
	rows, err := bucketDefault.ExecuteN1qlQuery(query, []interface{}{name, firmware})
	var data map[string]interface{}
	if err != nil {
		return data, err
	}

	if rows.Next(&data) {
		return data, nil
	}

	return data, errors.New("Function not found " + name)
}

func produceKafkaFunctionMessage(req *pb.CallDeviceFunctionRequest, functionDetails map[string]interface{}, firmware int32) error {
	msgValue, err := buildAvroMessage(functionDetails, req, firmware)
	if err != nil {
		return err
	}
	err = avroProducer.Produce(
		req.GetDeviceEUI(),
		msgValue,
		nil,
	)
	return err
}

func getAsBytes(value string, ptype interface{}) ([]byte, error) {
	if value == "" {
		return nil, errors.New("Missing required parameter")
	}

	buf := new(bytes.Buffer)
	if ptype == "i" {
		//4 byte signed int
		//ParseInt returns int64
		int64Value, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, err
		}
		int32Value := int32(int64Value)
		err = binary.Write(buf, binary.BigEndian, int32Value)
		if err != nil {
			return nil, err
		}

	} else if ptype == "t" {
		//2 byte signed int
		//ParseInt returns int64 but you can specify that it should fit into int16
		int64Value, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return nil, err
		}
		int16Value := int16(int64Value)
		err = binary.Write(buf, binary.BigEndian, int16Value)
		if err != nil {
			return nil, err
		}
	} else if ptype == "b" {
		//2 byte bool
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		//convert to 2 byte int
		var int16Value uint16
		int16Value = 0
		if boolValue {
			int16Value = 1
		}
		err = binary.Write(buf, binary.BigEndian, int16Value)
		if err != nil {
			return nil, err
		}
	} else {
		//string. ptype is length of string
		val := ptype.(float64)
		strlength := int(val)

		buf.Write([]byte(value))
		//check length
		if buf.Len() > strlength {
			return nil, errors.New("String too long, length: " +
				strconv.Itoa(buf.Len()) + ", allowed: " + strconv.Itoa(strlength))
		}
		//padding
		diff := strlength - buf.Len()
		for i := 0; i < diff; i++ {
			buf.WriteByte(byte(0))
		}
	}

	return buf.Bytes(), nil
}

/*
func getEmptyBytes(ptype interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if ptype == "i" {
		//4 byte signed int
		int32Value := int32(0)
		err := binary.Write(buf, binary.BigEndian, int32Value)
		if err != nil {
			return nil, err
		}

	} else if ptype == "t" || ptype == "b" {
		//2 byte signed int
		int16Value := int16(0)
		err := binary.Write(buf, binary.BigEndian, int16Value)
		if err != nil {
			return nil, err
		}
	} else {
		//string. ptype is length of string
		val := ptype.(float64)
		strlength := int(val)

		//padding
		for i := 0; i < strlength; i++ {
			buf.WriteByte(byte(0))
		}
	}

	return buf.Bytes(), nil
}*/

func buildAvroMessage(functionDetails map[string]interface{}, req *pb.CallDeviceFunctionRequest, firmwareVersion int32) (map[string]interface{}, error) {
	functionMsg := make(map[string]interface{})
	functionMsg["slot"] = 0
	functionMsg["index"] = functionDetails["index"]
	functionMsg["firmware"] = firmwareVersion

	if functionDetails["param1type"] != 0.0 {
		val, err := getAsBytes(req.GetParam1(), functionDetails["param1type"])
		if err != nil {
			return functionMsg, err
		}
		functionMsg["param1"] = val
	}

	if functionDetails["param2type"] != 0.0 {
		val, err := getAsBytes(req.GetParam2(), functionDetails["param2type"])
		if err != nil {
			return functionMsg, err
		}
		functionMsg["param2"] = val
	}

	if functionDetails["param3type"] != 0.0 {
		val, err := getAsBytes(req.GetParam3(), functionDetails["param3type"])
		if err != nil {
			return functionMsg, err
		}
		functionMsg["param3"] = val
	}

	if functionDetails["param4type"] != 0.0 {
		val, err := getAsBytes(req.GetParam4(), functionDetails["param4type"])
		if err != nil {
			return functionMsg, err
		}
		functionMsg["param4"] = val
	}
	fmt.Println(functionMsg)
	return functionMsg, nil
}

func checkToken(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Println("Could not get metadata")
		return errors.New("Could not get metadata")
	}
	token := md["authorization"]
	authRequest := &pbAuth.AuthRequest{
		Token:        token[0],
		AllowedRoles: []string{rolePowerpilotAdmin, rolePowerpilotInstaller},
	}
	authResponse, err := grpcAuthClient.CheckAuth(context.Background(), authRequest)
	if err != nil {
		return err
	}
	if authResponse.GetResult() == false {
		return errors.New(authResponse.GetMessage())
	}

	return nil
}

//gRPC
func (s *functionServiceServer) CallDeviceFunction(ctx context.Context, req *pb.CallDeviceFunctionRequest) (*pb.Response, error) {
	err := checkToken(ctx)
	if err != nil {
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	firmware, err := getFirmwareVersion(req.GetDeviceEUI())
	if err != nil {
		fmt.Println("firmware error")
		fmt.Println(err)
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}
	functionDetails, err := getFunctionDetailsByName(req.GetFunctionName(), firmware)
	if err != nil {
		fmt.Println("function details error")
		fmt.Println(err)
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	//produce kafka message
	err = produceKafkaFunctionMessage(req, functionDetails, firmware)
	if err != nil {
		fmt.Println(err)
		return &pb.Response{
			Reply: "NOT OK",
		}, err
	}

	fmt.Println("Send function complete")
	fmt.Println(req)

	//send response
	return &pb.Response{
		Reply: "OK",
	}, nil
}

func (s *functionServiceServer) GetDeviceFunctions(ctx context.Context, req *pb.GetDeviceFunctionsRequest) (*pb.DeviceFunctions, error) {
	err := checkToken(ctx)
	if err != nil {
		return nil, err
	}

	firmware, err := getFirmwareVersion(req.GetIdentifier())
	if err != nil {
		fmt.Println("firmware error")
		return nil, err
	}

	functions, err := getFunctionsForFirmware(firmware)
	if err != nil {
		fmt.Println("Error getting functions for firmware")
		fmt.Println(err)
		return nil, err
	}

	//send response
	return functions, nil
}

//end gRPC

func main() {
	conn, err := grpc.Dial(grpcAuthServerAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	grpcAuthClient = pbAuth.NewAuthServiceClient(conn)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDeviceFunctionServiceServer(grpcServer, &functionServiceServer{})
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}
