# Social Record FrontEnd  
The frontend for Social Record

---

## Table of content

  - [Table of content](#table-of-content)
  - [About](#about)
  - [Getting started](#getting-started)
  - [Application Design](#application-design)
  - [Usage](#usage)
  - [Testing](#testing)
  - [Collaboration](#collaboration)

## About

The Social Record FrontEnd enables active social media users to use Social Record and display the analyzed data from instagram of the user. The FrontEnd is created to show the data as part of an exhbition where the visitors can experience what public data social media providers expose about them. The goal is to to educate the visitors about data privacy and supports them to improve their abilities to control their data in the social web. 

Due to privacy reasons it is not planned to launch Social Record in the public web. To learn more about Social Record read the [general documentation](https://github.com/codeuniversity/smag-mvp).

## Getting started

**1. clone repo**

[https://github.com/codeuniversity/smag-mvp.git] (https://github.com/codeuniversity/smag-mvp.git)

**2. Switch to the FrontEnd folder**

`cd frontend`

**3. Install all the npm packages. Type the following command to install all npm packages**

`npm install`

**4. In order to run the application in development mode type the following command**

`npm start`

The application runs on localhost:3000

**5. When you’re ready to deploy to production you can create a minified bundle**

`npm run build`

**6. To connect to the Social Record database you need to connect the node.js server to the Social Record kubernetes cluster. You can also run the application locaon a local test environment.** 

1. Connect to the Kubernetes cluster of Social Records 
2. Run kubectl: ```kubectl get pods```
3. Select the pod ID from: ```envoy-proxy-deployment``` 
4. Forward the port from 4000 to 8000. Insert the pod ID instead of PODID: ```kubectl port-forward envoy-proxy-deployment-PODID 4000:8080```

## Application Design

**Components**
1. **App** component: This component enables the user to search for the username. 
2. **Result** component: This component displays the instagram-data (Username, Realname, Authorpicture, Bio, Posts). The component gets its data from the gRPC API. 
3. **Start** component: This component enables the user to take a picture of him/herself to identify the user using face recognition. 
4. **Title** component: This component displays the h1 headline in the App component and Start component.
5. **Form** component: This component displays the search form in the App component.
6. **IGPost** component: This component displays the instagram post in the Result component. 
7. **BackButton** component: This component displays a button in the Result component to go back the App component
8. **camera-feed** component: This component displays the camera functionality in the Start component. 

All styles are defined in the `index.css`

## Usage

To use the FrontEnd and access the gRPC API you have to import the auto-generated proto files in the respective components. 

```javascript
import {User, UserNameRequest, UserIdRequest, Post, UserIdResponse, UserSearchResponse} from "./protofiles/usersearch_pb.js";
import { UserSearchServiceClient } from "./protofiles/usersearch_grpc_web_pb";

const userSearch = new UserSearchServiceClient("http://localhost:4000");

const requestUser = new UserNameRequest();
    requestUser.setUserName(userName);

    userSearch.getUserWithUsername(requestUser, {}, (err, response) => {
    //...
});
```

The default address for the database is localhost. If you want to change that simply add the enviroment variable GRPC_POSTGRES_HOST to the grpc-servercontainer. 


## Testing
To test Social Record on a local machine you can follow the guideline in the [general documentation](https://github.com/codeuniversity/smag-mvp).

## Collaboration



