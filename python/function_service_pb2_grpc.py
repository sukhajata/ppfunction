# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import function_service_pb2 as function__service__pb2


class DeviceFunctionServiceStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.CallDeviceFunction = channel.unary_unary(
        '/ppfunction.DeviceFunctionService/CallDeviceFunction',
        request_serializer=function__service__pb2.CallDeviceFunctionRequest.SerializeToString,
        response_deserializer=function__service__pb2.Response.FromString,
        )
    self.GetDeviceFunctions = channel.unary_unary(
        '/ppfunction.DeviceFunctionService/GetDeviceFunctions',
        request_serializer=function__service__pb2.GetDeviceFunctionsRequest.SerializeToString,
        response_deserializer=function__service__pb2.DeviceFunctions.FromString,
        )


class DeviceFunctionServiceServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def CallDeviceFunction(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetDeviceFunctions(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_DeviceFunctionServiceServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'CallDeviceFunction': grpc.unary_unary_rpc_method_handler(
          servicer.CallDeviceFunction,
          request_deserializer=function__service__pb2.CallDeviceFunctionRequest.FromString,
          response_serializer=function__service__pb2.Response.SerializeToString,
      ),
      'GetDeviceFunctions': grpc.unary_unary_rpc_method_handler(
          servicer.GetDeviceFunctions,
          request_deserializer=function__service__pb2.GetDeviceFunctionsRequest.FromString,
          response_serializer=function__service__pb2.DeviceFunctions.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'ppfunction.DeviceFunctionService', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
