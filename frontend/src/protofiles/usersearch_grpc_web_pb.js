/* eslint-disable */
/**
 * @fileoverview gRPC-Web generated client stub for proto
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!

const grpc = {};
grpc.web = require("grpc-web");

const proto = {};
proto.proto = require("./usersearch_pb.js");

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.proto.UserSearchServiceClient = function(hostname, credentials, options) {
  if (!options) options = {};
  options["format"] = "text";

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
proto.proto.UserSearchServicePromiseClient = function(
  hostname,
  credentials,
  options
) {
  if (!options) options = {};
  options["format"] = "text";

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
 *   !proto.proto.UserNameRequest,
 *   !proto.proto.UserSearchResponse>}
 */
const methodDescriptor_UserSearchService_GetAllUsersLikeUsername = new grpc.web.MethodDescriptor(
  "/proto.UserSearchService/GetAllUsersLikeUsername",
  grpc.web.MethodType.UNARY,
  proto.proto.UserNameRequest,
  proto.proto.UserSearchResponse,
  /**
   * @param {!proto.proto.UserNameRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.UserSearchResponse.deserializeBinary
);

/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.UserNameRequest,
 *   !proto.proto.UserSearchResponse>}
 */
const methodInfo_UserSearchService_GetAllUsersLikeUsername = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.UserSearchResponse,
  /**
   * @param {!proto.proto.UserNameRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.UserSearchResponse.deserializeBinary
);

/**
 * @param {!proto.proto.UserNameRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.UserSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.UserSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.UserSearchServiceClient.prototype.getAllUsersLikeUsername = function(
  request,
  metadata,
  callback
) {
  return this.client_.rpcCall(
    this.hostname_ + "/proto.UserSearchService/GetAllUsersLikeUsername",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_GetAllUsersLikeUsername,
    callback
  );
};

/**
 * @param {!proto.proto.UserNameRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.UserSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.proto.UserSearchServicePromiseClient.prototype.getAllUsersLikeUsername = function(
  request,
  metadata
) {
  return this.client_.unaryCall(
    this.hostname_ + "/proto.UserSearchService/GetAllUsersLikeUsername",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_GetAllUsersLikeUsername
  );
};

/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.UserIdRequest,
 *   !proto.proto.User>}
 */
const methodDescriptor_UserSearchService_GetUserWithUserId = new grpc.web.MethodDescriptor(
  "/proto.UserSearchService/GetUserWithUserId",
  grpc.web.MethodType.UNARY,
  proto.proto.UserIdRequest,
  proto.proto.User,
  /**
   * @param {!proto.proto.UserIdRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.User.deserializeBinary
);

/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.UserIdRequest,
 *   !proto.proto.User>}
 */
const methodInfo_UserSearchService_GetUserWithUserId = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.User,
  /**
   * @param {!proto.proto.UserIdRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.User.deserializeBinary
);

/**
 * @param {!proto.proto.UserIdRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.UserSearchServiceClient.prototype.getUserWithUserId = function(
  request,
  metadata,
  callback
) {
  return this.client_.rpcCall(
    this.hostname_ + "/proto.UserSearchService/GetUserWithUserId",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_GetUserWithUserId,
    callback
  );
};

/**
 * @param {!proto.proto.UserIdRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.User>}
 *     A native promise that resolves to the response
 */
proto.proto.UserSearchServicePromiseClient.prototype.getUserWithUserId = function(
  request,
  metadata
) {
  return this.client_.unaryCall(
    this.hostname_ + "/proto.UserSearchService/GetUserWithUserId",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_GetUserWithUserId
  );
};

/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.UserNameRequest,
 *   !proto.proto.User>}
 */
const methodDescriptor_UserSearchService_GetUserWithUsername = new grpc.web.MethodDescriptor(
  "/proto.UserSearchService/GetUserWithUsername",
  grpc.web.MethodType.UNARY,
  proto.proto.UserNameRequest,
  proto.proto.User,
  /**
   * @param {!proto.proto.UserNameRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.User.deserializeBinary
);

/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.UserNameRequest,
 *   !proto.proto.User>}
 */
const methodInfo_UserSearchService_GetUserWithUsername = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.User,
  /**
   * @param {!proto.proto.UserNameRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.User.deserializeBinary
);

/**
 * @param {!proto.proto.UserNameRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.UserSearchServiceClient.prototype.getUserWithUsername = function(
  request,
  metadata,
  callback
) {
  return this.client_.rpcCall(
    this.hostname_ + "/proto.UserSearchService/GetUserWithUsername",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_GetUserWithUsername,
    callback
  );
};

/**
 * @param {!proto.proto.UserNameRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.User>}
 *     A native promise that resolves to the response
 */
proto.proto.UserSearchServicePromiseClient.prototype.getUserWithUsername = function(
  request,
  metadata
) {
  return this.client_.unaryCall(
    this.hostname_ + "/proto.UserSearchService/GetUserWithUsername",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_GetUserWithUsername
  );
};

/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.UserIdRequest,
 *   !proto.proto.InstaPostsResponse>}
 */
const methodDescriptor_UserSearchService_GetInstaPostsWithUserId = new grpc.web.MethodDescriptor(
  "/proto.UserSearchService/GetInstaPostsWithUserId",
  grpc.web.MethodType.UNARY,
  proto.proto.UserIdRequest,
  proto.proto.InstaPostsResponse,
  /**
   * @param {!proto.proto.UserIdRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.InstaPostsResponse.deserializeBinary
);

/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.UserIdRequest,
 *   !proto.proto.InstaPostsResponse>}
 */
const methodInfo_UserSearchService_GetInstaPostsWithUserId = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.InstaPostsResponse,
  /**
   * @param {!proto.proto.UserIdRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.InstaPostsResponse.deserializeBinary
);

/**
 * @param {!proto.proto.UserIdRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.InstaPostsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.InstaPostsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.UserSearchServiceClient.prototype.getInstaPostsWithUserId = function(
  request,
  metadata,
  callback
) {
  return this.client_.rpcCall(
    this.hostname_ + "/proto.UserSearchService/GetInstaPostsWithUserId",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_GetInstaPostsWithUserId,
    callback
  );
};

/**
 * @param {!proto.proto.UserIdRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.InstaPostsResponse>}
 *     A native promise that resolves to the response
 */
proto.proto.UserSearchServicePromiseClient.prototype.getInstaPostsWithUserId = function(
  request,
  metadata
) {
  return this.client_.unaryCall(
    this.hostname_ + "/proto.UserSearchService/GetInstaPostsWithUserId",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_GetInstaPostsWithUserId
  );
};

/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.FaceSearchRequest,
 *   !proto.proto.FaceSearchResponse>}
 */
const methodDescriptor_UserSearchService_SearchSimilarFaces = new grpc.web.MethodDescriptor(
  "/proto.UserSearchService/SearchSimilarFaces",
  grpc.web.MethodType.UNARY,
  proto.proto.FaceSearchRequest,
  proto.proto.FaceSearchResponse,
  /**
   * @param {!proto.proto.FaceSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.FaceSearchResponse.deserializeBinary
);

/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.FaceSearchRequest,
 *   !proto.proto.FaceSearchResponse>}
 */
const methodInfo_UserSearchService_SearchSimilarFaces = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.FaceSearchResponse,
  /**
   * @param {!proto.proto.FaceSearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.FaceSearchResponse.deserializeBinary
);

/**
 * @param {!proto.proto.FaceSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.FaceSearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.FaceSearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.UserSearchServiceClient.prototype.searchSimilarFaces = function(
  request,
  metadata,
  callback
) {
  return this.client_.rpcCall(
    this.hostname_ + "/proto.UserSearchService/SearchSimilarFaces",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_SearchSimilarFaces,
    callback
  );
};

/**
 * @param {!proto.proto.FaceSearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.FaceSearchResponse>}
 *     A native promise that resolves to the response
 */
proto.proto.UserSearchServicePromiseClient.prototype.searchSimilarFaces = function(
  request,
  metadata
) {
  return this.client_.unaryCall(
    this.hostname_ + "/proto.UserSearchService/SearchSimilarFaces",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_SearchSimilarFaces
  );
};

/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.WeightedPosts,
 *   !proto.proto.WeightedUsers>}
 */
const methodDescriptor_UserSearchService_SearchUsersWithWeightedPosts = new grpc.web.MethodDescriptor(
  "/proto.UserSearchService/SearchUsersWithWeightedPosts",
  grpc.web.MethodType.UNARY,
  proto.proto.WeightedPosts,
  proto.proto.WeightedUsers,
  /**
   * @param {!proto.proto.WeightedPosts} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.WeightedUsers.deserializeBinary
);

/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.WeightedPosts,
 *   !proto.proto.WeightedUsers>}
 */
const methodInfo_UserSearchService_SearchUsersWithWeightedPosts = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.WeightedUsers,
  /**
   * @param {!proto.proto.WeightedPosts} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.WeightedUsers.deserializeBinary
);

/**
 * @param {!proto.proto.WeightedPosts} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.WeightedUsers)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.WeightedUsers>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.UserSearchServiceClient.prototype.searchUsersWithWeightedPosts = function(
  request,
  metadata,
  callback
) {
  return this.client_.rpcCall(
    this.hostname_ + "/proto.UserSearchService/SearchUsersWithWeightedPosts",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_SearchUsersWithWeightedPosts,
    callback
  );
};

/**
 * @param {!proto.proto.WeightedPosts} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.WeightedUsers>}
 *     A native promise that resolves to the response
 */
proto.proto.UserSearchServicePromiseClient.prototype.searchUsersWithWeightedPosts = function(
  request,
  metadata
) {
  return this.client_.unaryCall(
    this.hostname_ + "/proto.UserSearchService/SearchUsersWithWeightedPosts",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_SearchUsersWithWeightedPosts
  );
};

/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.UserIdRequest,
 *   !proto.proto.UserDataPointCount>}
 */
const methodDescriptor_UserSearchService_DataPointCountForUserId = new grpc.web.MethodDescriptor(
  "/proto.UserSearchService/DataPointCountForUserId",
  grpc.web.MethodType.UNARY,
  proto.proto.UserIdRequest,
  proto.proto.UserDataPointCount,
  /**
   * @param {!proto.proto.UserIdRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.UserDataPointCount.deserializeBinary
);

/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.UserIdRequest,
 *   !proto.proto.UserDataPointCount>}
 */
const methodInfo_UserSearchService_DataPointCountForUserId = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.UserDataPointCount,
  /**
   * @param {!proto.proto.UserIdRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.UserDataPointCount.deserializeBinary
);

/**
 * @param {!proto.proto.UserIdRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.UserDataPointCount)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.UserDataPointCount>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.UserSearchServiceClient.prototype.dataPointCountForUserId = function(
  request,
  metadata,
  callback
) {
  return this.client_.rpcCall(
    this.hostname_ + "/proto.UserSearchService/DataPointCountForUserId",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_DataPointCountForUserId,
    callback
  );
};

/**
 * @param {!proto.proto.UserIdRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.UserDataPointCount>}
 *     A native promise that resolves to the response
 */
proto.proto.UserSearchServicePromiseClient.prototype.dataPointCountForUserId = function(
  request,
  metadata
) {
  return this.client_.unaryCall(
    this.hostname_ + "/proto.UserSearchService/DataPointCountForUserId",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_DataPointCountForUserId
  );
};

/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.UserIdRequest,
 *   !proto.proto.FoundCities>}
 */
const methodDescriptor_UserSearchService_FindCitiesForUserId = new grpc.web.MethodDescriptor(
  "/proto.UserSearchService/FindCitiesForUserId",
  grpc.web.MethodType.UNARY,
  proto.proto.UserIdRequest,
  proto.proto.FoundCities,
  /**
   * @param {!proto.proto.UserIdRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.FoundCities.deserializeBinary
);

/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.UserIdRequest,
 *   !proto.proto.FoundCities>}
 */
const methodInfo_UserSearchService_FindCitiesForUserId = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.FoundCities,
  /**
   * @param {!proto.proto.UserIdRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.FoundCities.deserializeBinary
);

/**
 * @param {!proto.proto.UserIdRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.FoundCities)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.FoundCities>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.UserSearchServiceClient.prototype.findCitiesForUserId = function(
  request,
  metadata,
  callback
) {
  return this.client_.rpcCall(
    this.hostname_ + "/proto.UserSearchService/FindCitiesForUserId",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_FindCitiesForUserId,
    callback
  );
};

/**
 * @param {!proto.proto.UserIdRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.FoundCities>}
 *     A native promise that resolves to the response
 */
proto.proto.UserSearchServicePromiseClient.prototype.findCitiesForUserId = function(
  request,
  metadata
) {
  return this.client_.unaryCall(
    this.hostname_ + "/proto.UserSearchService/FindCitiesForUserId",
    request,
    metadata || {},
    methodDescriptor_UserSearchService_FindCitiesForUserId
  );
};

module.exports = proto.proto;
