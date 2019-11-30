# SMAG gRPC Web API

## About

In our project we are using a [gRPC Web](https://grpc.io/docs/) API. For that we are using an [envoy proxy](https://www.envoyproxy.io/docs/envoy/latest/) to be able to connect to the gRPC Server. As our system is not publicly accessible an AWS Account in our Organisation with the appropriate access is required.

## Requirements

In order to successfully use our api make sure to have:

- a running [kubernetes setup](https://github.com/codeuniversity/smag-deploy/blob/master/README.md) (permssion required)
- _optional for local testing_: [protoc](http://google.github.io/proto-lens/installing-protoc.html) to generate the protofiles for the frontend

## Usage

To use the production enviroment do the following steps:

1. Get name of envoy proxy `kubectl get pods | grep envoy`
2. Forward the envoy-pod port with `kubectl port-forward envoy-proxy-deployment-6b89675d5b-d86c4 4000:8080`
3. To make use of the API in the React Frontend import and run the following:
   ```javascript
   import {
     User,
     UserNameRequest,
     UserIdRequest,
     InstaPostsResponse,
     UserSearchResponse
   } from "./protofiles/client/usersearch_pb.js";
   import { UserSearchServiceClient } from "./protofiles/client/usersearch_grpc_web_pb.js";
   var userSearch = new UserSearchServiceClient("http://localhost:4000");
   var request = new UserName();
   request.setUserName("codeuniversity");
   userSearch.getUserWithUsername(request, {}, function(err, response) {
     //example function call...
   });
   ```

## Functions

To check the attributes of the proto messages take a look at the protofile [userserach.proto](https://github.com/codeuniversity/smag-mvp/blob/master/api/proto/usersearch.proto)

| **Method** | **Function Name**        | **Input Message** | **Return Message** |
| ---------- | ------------------------ | ----------------- | ------------------ |
| GET        | getUserWithUsername      | UserNameRequest   | User               |
| GET        | getAllUsersLikeUsername  | UserNameRequest   | UserSearchResponse |
| GET        | getTaggedPostsWithUserId | UserIdRequest     | InstaPostsResponse |
| GET        | getInstaPostssWithUserid | UserIdRequest     | UserSearchResponse |
