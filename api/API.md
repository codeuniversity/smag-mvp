# gRPC API
## About
In our project we are using a [gRPC web](https://grpc.io/docs/) API. For that we are using an [envoy proxy](https://www.envoyproxy.io/docs/envoy/latest/) to be able to connect to the gRPC Server. As we have no real users in our Project and we don't want our system publicly accessible from the outside you need to have a AWS Account in our Organisation with the appropriate Access. The only way to access our API is through port-forwarding from the kubernetes cluster. 

## Requirements
- In order to successfully use our api make sure to have:
  * a AWS Account setup in our Organisation
  - [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/), to forward the ports you need for the to connect with the frontend 
  - [aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv1.html), to get your credentials to access the cluster
  - *optional for local testing*: [protoc](http://google.github.io/proto-lens/installing-protoc.html](http://google.github.io/proto-lens/installing-protoc.html) to generate the protofiles for the frontend

## Usage
to use the production enviroment do the following steps:
1. Get name of envoy proxy `kubectl get pods | grep envoy`
2. To be able to connect to the API you need to forward the envoy-pod port with  ```kubectl port-forward envoy-proxy-deployment-6b89675d5b-d86c4 4000:8080```
3. To make use of the API in the React Frontend import the include:
     ```javascript
     import {User, UserNameRequest, UserIdRequest, InstaPostsResponse, UserSearchResponse} from "./protofiles/client/usersearch_pb.js";
     import {UserSearchServiceClient} from "./protofiles/client/usersearch_grpc_web_pb.js";

     var userSearch = new UserSearchServiceClient('http://localhost:8080');
    ```
## Functions
to check the attributes of the proto messages take a look at the protofile *userserach.proto*

| **Method** | **function name**           | **input message**   | **return message** |
|------------|-----------------------------|---------------------|--------------------|
| GET        | getUserWithUsername         | UserNameRequest     | User               |
| GET        | getAllUsersWithUsername     | UserNameRequest     | UserSearchResponse |
| GET        | getTaggedPostsWithUserId    | UserIdRequest       | InstaPostsResponse |
| GET        | getInstaPostssWithUserid    | UserIdRequest       | UserSearchResponse |

## local testing
In the case of local changes to the API, it is possible to run the setup locally
1. Optional: re-generate the protofiles with `make gen-client`
2. `docker-compose up envoy-proxy`
3. initialize the Database with `make init-db`
4. follow the instructions for the frontend
5. Then connect with the envoy-proxy via `localhost:4000`
