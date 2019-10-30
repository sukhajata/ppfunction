/**
 * @fileoverview gRPC-Web generated client stub for ppfunction
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.ppfunction = require('./function-service_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.ppfunction.DeviceFunctionServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.ppfunction.DeviceFunctionServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ppfunction.CallDeviceFunctionRequest,
 *   !proto.ppfunction.Response>}
 */
const methodDescriptor_DeviceFunctionService_CallDeviceFunction = new grpc.web.MethodDescriptor(
  '/ppfunction.DeviceFunctionService/CallDeviceFunction',
  grpc.web.MethodType.UNARY,
  proto.ppfunction.CallDeviceFunctionRequest,
  proto.ppfunction.Response,
  /**
   * @param {!proto.ppfunction.CallDeviceFunctionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ppfunction.Response.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ppfunction.CallDeviceFunctionRequest,
 *   !proto.ppfunction.Response>}
 */
const methodInfo_DeviceFunctionService_CallDeviceFunction = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ppfunction.Response,
  /**
   * @param {!proto.ppfunction.CallDeviceFunctionRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ppfunction.Response.deserializeBinary
);


/**
 * @param {!proto.ppfunction.CallDeviceFunctionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ppfunction.Response)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ppfunction.Response>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ppfunction.DeviceFunctionServiceClient.prototype.callDeviceFunction =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ppfunction.DeviceFunctionService/CallDeviceFunction',
      request,
      metadata || {},
      methodDescriptor_DeviceFunctionService_CallDeviceFunction,
      callback);
};


/**
 * @param {!proto.ppfunction.CallDeviceFunctionRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ppfunction.Response>}
 *     A native promise that resolves to the response
 */
proto.ppfunction.DeviceFunctionServicePromiseClient.prototype.callDeviceFunction =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ppfunction.DeviceFunctionService/CallDeviceFunction',
      request,
      metadata || {},
      methodDescriptor_DeviceFunctionService_CallDeviceFunction);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.ppfunction.GetDeviceFunctionsRequest,
 *   !proto.ppfunction.DeviceFunctions>}
 */
const methodDescriptor_DeviceFunctionService_GetDeviceFunctions = new grpc.web.MethodDescriptor(
  '/ppfunction.DeviceFunctionService/GetDeviceFunctions',
  grpc.web.MethodType.UNARY,
  proto.ppfunction.GetDeviceFunctionsRequest,
  proto.ppfunction.DeviceFunctions,
  /**
   * @param {!proto.ppfunction.GetDeviceFunctionsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ppfunction.DeviceFunctions.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.ppfunction.GetDeviceFunctionsRequest,
 *   !proto.ppfunction.DeviceFunctions>}
 */
const methodInfo_DeviceFunctionService_GetDeviceFunctions = new grpc.web.AbstractClientBase.MethodInfo(
  proto.ppfunction.DeviceFunctions,
  /**
   * @param {!proto.ppfunction.GetDeviceFunctionsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.ppfunction.DeviceFunctions.deserializeBinary
);


/**
 * @param {!proto.ppfunction.GetDeviceFunctionsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.ppfunction.DeviceFunctions)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.ppfunction.DeviceFunctions>|undefined}
 *     The XHR Node Readable Stream
 */
proto.ppfunction.DeviceFunctionServiceClient.prototype.getDeviceFunctions =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/ppfunction.DeviceFunctionService/GetDeviceFunctions',
      request,
      metadata || {},
      methodDescriptor_DeviceFunctionService_GetDeviceFunctions,
      callback);
};


/**
 * @param {!proto.ppfunction.GetDeviceFunctionsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.ppfunction.DeviceFunctions>}
 *     A native promise that resolves to the response
 */
proto.ppfunction.DeviceFunctionServicePromiseClient.prototype.getDeviceFunctions =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/ppfunction.DeviceFunctionService/GetDeviceFunctions',
      request,
      metadata || {},
      methodDescriptor_DeviceFunctionService_GetDeviceFunctions);
};


module.exports = proto.ppfunction;

