
#commands for generating client server code with protoc

#go
export PATH=$PATH:~/go/bin

protoc --go_out=plugins=grpc:. *.proto 

go install

#web gRPC / js
protoc -I=. *.proto \
  --js_out=import_style=commonjs:web \
  --grpc-web_out=import_style=commonjs,mode=grpcwebtext:web

#add /* eslint-disable */ to config-service.pb.go
#add dependencies to package.json
  
#python
pip install grpcio-tools

python -m grpc_tools.protoc -I . --python_out=python --grpc_python_out=python *.proto