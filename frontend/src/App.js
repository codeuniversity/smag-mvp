import React, { Component } from "react";
import Form from "./components/Form";
import Title from "./components/Title";
import Results from "./components/Results";
import "./index.css";
import {
  User,
  UserSearchRequest,
  UserSearchResponse
} from "./protofiles/usersearch_pb.js";
import { UserSearchServiceClient } from "./protofiles/usersearch_grpc_web_pb";
import { withRouter } from "react-router";
import PropTypes from "prop-types";

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
      this.props.history.push({
        pathname: "/results",
        state: { results: userdata }
      });

      //this.setState({ users: userdata });
    });
  };
  render() {
    return (
      <div className="container-center">
        <div className="column-center">
          <Title />
          <Form onSubmit={this.handleSubmit} />
        </div>
      </div>
    );
  }
}

App.propTypes = {
  history: PropTypes.shape({
    push: PropTypes.shape({
      pathname: PropTypes.string,
      state: PropTypes.object
    })
  })
};

export default withRouter(App);
