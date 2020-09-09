/**
 * @fileoverview gRPC-Web generated client stub for proto
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')
const proto = {};
proto.proto = require('./server_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.proto.ViewClient =
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
proto.proto.ViewPromiseClient =
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
 *   !proto.proto.SearchRequest,
 *   !proto.proto.SearchResponse>}
 */
const methodDescriptor_View_SearchQuery = new grpc.web.MethodDescriptor(
  '/proto.View/SearchQuery',
  grpc.web.MethodType.UNARY,
  proto.proto.SearchRequest,
  proto.proto.SearchResponse,
  /**
   * @param {!proto.proto.SearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.SearchResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.SearchRequest,
 *   !proto.proto.SearchResponse>}
 */
const methodInfo_View_SearchQuery = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.SearchResponse,
  /**
   * @param {!proto.proto.SearchRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.SearchResponse.deserializeBinary
);


/**
 * @param {!proto.proto.SearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.SearchResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.SearchResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.ViewClient.prototype.searchQuery =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/proto.View/SearchQuery',
      request,
      metadata || {},
      methodDescriptor_View_SearchQuery,
      callback);
};


/**
 * @param {!proto.proto.SearchRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.SearchResponse>}
 *     A native promise that resolves to the response
 */
proto.proto.ViewPromiseClient.prototype.searchQuery =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/proto.View/SearchQuery',
      request,
      metadata || {},
      methodDescriptor_View_SearchQuery);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.RankingsRequest,
 *   !proto.proto.RankingsResponse>}
 */
const methodDescriptor_View_ViewRankings = new grpc.web.MethodDescriptor(
  '/proto.View/ViewRankings',
  grpc.web.MethodType.UNARY,
  proto.proto.RankingsRequest,
  proto.proto.RankingsResponse,
  /**
   * @param {!proto.proto.RankingsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.RankingsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.RankingsRequest,
 *   !proto.proto.RankingsResponse>}
 */
const methodInfo_View_ViewRankings = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.RankingsResponse,
  /**
   * @param {!proto.proto.RankingsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.RankingsResponse.deserializeBinary
);


/**
 * @param {!proto.proto.RankingsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.RankingsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.RankingsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.ViewClient.prototype.viewRankings =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/proto.View/ViewRankings',
      request,
      metadata || {},
      methodDescriptor_View_ViewRankings,
      callback);
};


/**
 * @param {!proto.proto.RankingsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.RankingsResponse>}
 *     A native promise that resolves to the response
 */
proto.proto.ViewPromiseClient.prototype.viewRankings =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/proto.View/ViewRankings',
      request,
      metadata || {},
      methodDescriptor_View_ViewRankings);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.YrTotalRequest,
 *   !proto.proto.YrTotalResponse>}
 */
const methodDescriptor_View_ViewYrTotals = new grpc.web.MethodDescriptor(
  '/proto.View/ViewYrTotals',
  grpc.web.MethodType.UNARY,
  proto.proto.YrTotalRequest,
  proto.proto.YrTotalResponse,
  /**
   * @param {!proto.proto.YrTotalRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.YrTotalResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.YrTotalRequest,
 *   !proto.proto.YrTotalResponse>}
 */
const methodInfo_View_ViewYrTotals = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.YrTotalResponse,
  /**
   * @param {!proto.proto.YrTotalRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.YrTotalResponse.deserializeBinary
);


/**
 * @param {!proto.proto.YrTotalRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.YrTotalResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.YrTotalResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.ViewClient.prototype.viewYrTotals =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/proto.View/ViewYrTotals',
      request,
      metadata || {},
      methodDescriptor_View_ViewYrTotals,
      callback);
};


/**
 * @param {!proto.proto.YrTotalRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.YrTotalResponse>}
 *     A native promise that resolves to the response
 */
proto.proto.ViewPromiseClient.prototype.viewYrTotals =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/proto.View/ViewYrTotals',
      request,
      metadata || {},
      methodDescriptor_View_ViewYrTotals);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.GetObjRequest,
 *   !proto.proto.GetObjResponse>}
 */
const methodDescriptor_View_ViewObject = new grpc.web.MethodDescriptor(
  '/proto.View/ViewObject',
  grpc.web.MethodType.UNARY,
  proto.proto.GetObjRequest,
  proto.proto.GetObjResponse,
  /**
   * @param {!proto.proto.GetObjRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.GetObjResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.GetObjRequest,
 *   !proto.proto.GetObjResponse>}
 */
const methodInfo_View_ViewObject = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.GetObjResponse,
  /**
   * @param {!proto.proto.GetObjRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.GetObjResponse.deserializeBinary
);


/**
 * @param {!proto.proto.GetObjRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.GetObjResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.GetObjResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.ViewClient.prototype.viewObject =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/proto.View/ViewObject',
      request,
      metadata || {},
      methodDescriptor_View_ViewObject,
      callback);
};


/**
 * @param {!proto.proto.GetObjRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.GetObjResponse>}
 *     A native promise that resolves to the response
 */
proto.proto.ViewPromiseClient.prototype.viewObject =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/proto.View/ViewObject',
      request,
      metadata || {},
      methodDescriptor_View_ViewObject);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.proto.Empty,
 *   !proto.proto.Empty>}
 */
const methodDescriptor_View_NoOp = new grpc.web.MethodDescriptor(
  '/proto.View/NoOp',
  grpc.web.MethodType.UNARY,
  proto.proto.Empty,
  proto.proto.Empty,
  /**
   * @param {!proto.proto.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.Empty.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.Empty,
 *   !proto.proto.Empty>}
 */
const methodInfo_View_NoOp = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.Empty,
  /**
   * @param {!proto.proto.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.Empty.deserializeBinary
);


/**
 * @param {!proto.proto.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.ViewClient.prototype.noOp =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/proto.View/NoOp',
      request,
      metadata || {},
      methodDescriptor_View_NoOp,
      callback);
};


/**
 * @param {!proto.proto.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.Empty>}
 *     A native promise that resolves to the response
 */
proto.proto.ViewPromiseClient.prototype.noOp =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/proto.View/NoOp',
      request,
      metadata || {},
      methodDescriptor_View_NoOp);
};


module.exports = proto.proto;

