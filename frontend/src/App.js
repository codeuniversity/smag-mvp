import React, { Component } from 'react';
import Form from './Form.js';
import Title from './Title.js';
import './index.css';
import {User, UserSearchRequest, UserSearchResponse} from "./protofiles/usersearch_pb.js";
import {UserSearchServiceClient} from "./protofiles/usersearch_grpc_web_pb";
import {withRouter} from 'react-router';
// import {saveUser} from './Api';
// eslint-disable-next-line

  class App extends Component {
    
    componentDidMount() {
 

  }

    handleSubmit = (userName) => {
      // saveUser(user).then(() =>
      //   this.props.history.push('/dashboard')
      // )

      const userSearch = new UserSearchServiceClient('http://localhost:4000');
      const request = new UserSearchRequest();
      console.log(request);
      // request.setUserName("lukasmenzel");
      request.setUserName(userName);
      userSearch.getAllUsersLikeUsername(request, {},function(err, response) {
        console.log(err)
        console.log(response.getUserListList())
      });
    }

    render() {
      return (
        <div className="container">
        <div className="column">
          <Title />
          <Form onSubmit={this.handleSubmit}/>              
          </div>
      </div>
      )
    }
  }

export default withRouter(App);
