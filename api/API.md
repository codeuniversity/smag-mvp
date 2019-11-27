# gRPC API

- [Requirements](#Requirements)
- [Usage](#usage)
- [Functions](#functions)
- [Testing](#testing)

## Requirements
- In order to successfully use our api make sure to have:
    - [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed
    - [aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv1.html) installed
    - a AWS Account setup in our Organisation

## Usage
- Make sure to `npm install google-protobuf grpc-web`
- To be able to connect to the API you also need ot forward the envoy-pod port with  `kubectl port-forward envoy-proxy-deployment-6b89675d5b-d86c4 4000:8080`
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
|method|function name                                                |Description|
|------|-------------------------------------------------------------|-----------|
|GET   |getUserWithUsername(UserNameRequest) User                    |
|GET   |getAllUsersWithUsername(UserNameRequest) UserSearchResponse  |
|GET   |getInstaPostssWithUserid(UserIdRequest) InstaPostsResponse   |
|GET   |getTaggedPostsWithUserId(UserIdRequest) InstaPostsResponse   |

- `getUserWithUsername(UserNameRequest) User`
    > Queries the Database for one specific User
- `getAllUsersWithUsername(UserNameRequest) UserSearchResponse`
    > Queries the database for all users that have a similar usenames and returns array of user
- `getInstaPostssWithUserid(UserIdRequest) InstaPostsResponse`
    > GetInstaPostsWithUserId takes the User id and returns all Instagram Posts of a User
- `getTaggedPostsWithUserId(UserIdRequest) InstaPostsResponse`
    > GetTaggedPostsWithUserId returns all Posts the given User is tagged on

## Testing
1. `docker-compose up`
1. initialize the Database with `make init-db`
1. Then connect with the envoy-proxy via `localhost:4000`
