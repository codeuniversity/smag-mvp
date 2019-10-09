# gRPC API

## Usage
- Make sure to `npm install google-protobuf grpc-web`
- Then import the auto-generated proto files
    ```javascript
    import {User, UserName, UserSearchResponse} from "./proto/client/usersearch_pb.js";
    import {UserSearchServiceClient} from "./proto/client/usersearch_grpc_web_pb.js";

    var userSearch = new UserSearchServiceClient('http://localhost:8080');

    var request = new UserName();
    request.setUserName("codeuniversity");
    
    userSearch.getUserWithUsername(request, {},function(err, response) {
        //...
    });
    ```
## Functions
- `getUserWithUsername(UserSearchRequest) User`
    - Queries the Database for one specific User
- `getAllUsersWithUsername(UserSearchRequest) UserSearchRsponse`
    - Queries the database for all users that have a similar usenames and returns array of user


## Testing

1. Start postgres container
1. Change the IP Adress in the envoy.yaml config to yours (WIP)
    ```yaml
    hosts: [{ socket_address: { address: <your ip>, port_value: 10000 }}]
    ```
1.start gRPC Server with `go run api/grpcserver/main/main.go`
1. build and start the envoy Container
    `docker build -t envoy-proxy -f api/envoy-proxy/Dockerfile .`
    `docker run -p 8080:8080 envoy-proxy`


