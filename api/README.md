# gRPC API

- [Usage](#usage)
- [Functions](#functions)
- [Testing](#testing)

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
- The default address for the database is `localhost`.
  If you want to change that simply add the enviroment variable `GRPC_POSTGRES_HOST` to the `grpc-server`container


## Functions
- `getUserWithUsername(UserNameRequest) User`
    > Queries the Database for one specific User
- `getAllUsersWithUsername(UserNameRequest) UserSearchResponse`
    > Queries the database for all users that have a similar usenames and returns array of user
- `getInstaPostssWithUserid(UserIdRequest) InstaPostsResponse`
    > GetInstaPostsWithUserId takes the User id and returns all Instagram Posts of a User

## Testing
1. `docker-compose up`
1. initialize the Database with `make init-db`
1. Then connect with the envoy-proxy via `localhost:4000`
