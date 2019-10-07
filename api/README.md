# gRPC API

## Usage

```javascript
import {User, UserName, UserSearchResponse} from "./proto/client/usersearch_pb.js";
import {UserSearchServiceClient} from "./proto/client/usersearch_grpc_web_pb";

var userSearch = new UserSearchServiceClient('http://localhost:8080');

var request = new UserName();
request.setUserName("codeuniversity");
  
userSearch.getUserWithUsername(request, {},function(err, response) {
    //...
});
```

## Testing

- for testing you first of all have to start the envoy Docker container 


