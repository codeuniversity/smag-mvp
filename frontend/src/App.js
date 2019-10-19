import React, { Component } from "react";
import Form from "./Form.js";
import Title from "./Title.js";
import "./index.css";
import {
  User,
  UserSearchRequest,
  UserSearchResponse
} from "./protofiles/usersearch_pb.js";
import { UserSearchServiceClient } from "./protofiles/usersearch_grpc_web_pb";
import { withRouter } from "react-router";

// eslint-disable-next-line

class App extends Component {
  state = {
    users: []
  };

  handleSubmit = userName => {
    const userSearch = new UserSearchServiceClient("http://localhost:4000");
    const request = new UserSearchRequest();

    request.setUserName(userName);
    userSearch.getAllUsersLikeUsername(request, {}, (err, response) => {
      console.log(err);
      const users = response.getUserListList();
      users.map(user => ({
        bio: user.getBio(),
        avatarurl: user.getAvatarUrl(),
        username: user.getUserName(),
        realname: user.getRealName()
      }));
      const userdata = users.map(user => ({
        bio: user.getBio(),
        avatarurl: user.getAvatarUrl(),
        username: user.getUserName(),
        realname: user.getRealName()
      }));

      this.setState({ users: userdata });
    });
  };

  render() {
    console.log(this.state.users);
    return (
      <div className="container">
        <div className="column">
          <Title />
          <Form onSubmit={this.handleSubmit} />
        </div>
      </div>
    );
  }
}

export default withRouter(App);
