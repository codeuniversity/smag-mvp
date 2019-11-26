#Social Record FrontEnd  
The frontend for Social Record

---

## Table of content

- [Table of content](#table-of-content)
- [About](#about)
- [Getting started](#getting-started)

## About

The Social Record FrontEnd enables active social media users to use Social Record and display the analyzed data from instagram of the user. The FrontEnd is created to show the data as part of an exhbition where the visitors can experience what public data social media providers expose about them. It is not planned to present Social Record in the public web due to privacy reasons. To learn more about Social Record read the [general documentation](https://github.com/codeuniversity/smag-mvp).

## Getting started

**clone repo**

`https://github.com/codeuniversity/smag-mvp.git`

Switch to the FrontEnd folder:
`cd frontend`

Install all the npm packages. Type the following command to install all npm packages:

`npm install`

In order to run the application type the following command:

`npm start`

The application runs on localhost:3000

To connect to the Social Record database you need to connect the node.js server to the Social Record kubernetes cluster. You can also run the application locaon a local test environment. 

1. Connect to the Kubernetes cluster of Social Records 
2. Run kubectl: ```kubectl get pods```
3. Select the pod ID from: ```envoy-proxy-deployment``` 
4. Forward the port from 4000 to 8000. Insert the pod ID instead of PODID: ```kubectl port-forward envoy-proxy-deployment-PODID 4000:8080```


##Functions


##Testing



##Collaboration

