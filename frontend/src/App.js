import React, { Component } from "react";
import Form from "./components/Form";
import H1 from "./components/H1";
import Result from "./components/Result";
import "./index.css";
import {
  User,
  UserNameRequest,
  UserIdRequest,
  Post,
  UserIdResponse,
  UserSearchResponse
} from "./protofiles/usersearch_pb.js";
import { UserSearchServiceClient } from "./protofiles/usersearch_grpc_web_pb";
import { withRouter } from "react-router";
import PropTypes from "prop-types";

// eslint-disable-next-line

class App extends Component {
  handleSubmit = userName => {
    const userSearch = new UserSearchServiceClient("http://localhost:4000");

    const requestUser = new UserNameRequest();

    requestUser.setUserName(userName);
    userSearch.getUserWithUsername(requestUser, {}, (err, response) => {
      if (err) {
        console.log(err);
        return;
      }
      const user = {
        id: response.getId(),
        bio: response.getBio(),
        avatarurl: response.getAvatarUrl(),
        username: response.getUserName(),
        realname: response.getRealName()
      };

      this.props.history.push({
        pathname: "/result",
        state: { user }
      });
    });
  };
  render() {
    return (
      <div className="container-center">
        <div className="column-center">
          <H1>Find out your public digital identity!</H1>
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
