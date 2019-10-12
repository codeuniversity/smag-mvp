/**
 * @fileoverview gRPC-Web generated client stub for proto
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.proto = require('./usersearch_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.proto.UserSearchServiceClient =
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

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.proto.UserSearchServicePromiseClient =
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

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.UserSearchRequest,
 *   !proto.proto.UserSearchResponse>}
 */
const methodDescriptor_UserSearchService_GetAllUsersLikeUsername = new grpc.web.MethodDescriptor(
  '/proto.UserSearchService/GetAllUsersLikeUsername',
  grpc.web.MethodType.UNARY,
  proto.proto.UserSearchRequest,
  proto.proto.UserSearchResponse,
  /** @param {!proto.proto.UserSearchRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.UserSearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.UserSearchRequest,
 *   !proto.proto.UserSearchResponse>}
 */
const methodInfo_UserSearchService_GetAllUsersLikeUsername = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.UserSearchResponse,
  /** @param {!proto.proto.UserSearchRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.UserSearchResponse.deserializeBinary
);


/**
 * @param {!proto.proto.UserSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.UserSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.UserSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.UserSearchServiceClient.prototype.getAllUsersLikeUsername =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/proto.UserSearchService/GetAllUsersLikeUsername',
      request,
      metadata || {},
      methodDescriptor_UserSearchService_GetAllUsersLikeUsername,
      callback);
};


/**
 * @param {!proto.proto.UserSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.UserSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.proto.UserSearchServicePromiseClient.prototype.getAllUsersLikeUsername =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/proto.UserSearchService/GetAllUsersLikeUsername',
      request,
      metadata || {},
      methodDescriptor_UserSearchService_GetAllUsersLikeUsername);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.UserSearchRequest,
 *   !proto.proto.User>}
 */
const methodDescriptor_UserSearchService_GetUserWithUsername = new grpc.web.MethodDescriptor(
  '/proto.UserSearchService/GetUserWithUsername',
  grpc.web.MethodType.UNARY,
  proto.proto.UserSearchRequest,
  proto.proto.User,
  /** @param {!proto.proto.UserSearchRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.User.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.UserSearchRequest,
 *   !proto.proto.User>}
 */
const methodInfo_UserSearchService_GetUserWithUsername = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.User,
  /** @param {!proto.proto.UserSearchRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.User.deserializeBinary
);


/**
 * @param {!proto.proto.UserSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.UserSearchServiceClient.prototype.getUserWithUsername =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/proto.UserSearchService/GetUserWithUsername',
      request,
      metadata || {},
      methodDescriptor_UserSearchService_GetUserWithUsername,
      callback);
};


/**
 * @param {!proto.proto.UserSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.User>}
 *     A native promise that resolves to the response
 */
proto.proto.UserSearchServicePromiseClient.prototype.getUserWithUsername =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/proto.UserSearchService/GetUserWithUsername',
      request,
      metadata || {},
      methodDescriptor_UserSearchService_GetUserWithUsername);
};


module.exports = proto.proto;

