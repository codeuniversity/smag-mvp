import React, { Component } from "react";
import Form from "../components/Form";
import H2 from "../components/H2";
import Button from "../components/Button";
import Result from "../components/Result";
import {
  User,
  UserNameRequest,
  UserIdRequest,
  Post,
  UserIdResponse,
  UserSearchResponse
} from "../protofiles/usersearch_pb.js";
import { UserSearchServiceClient } from "../protofiles/usersearch_grpc_web_pb";
import { withRouter } from "react-router";
import PropTypes from "prop-types";

// eslint-disable-next-line

class SearchProfile extends Component {
  handleSubmit = async userName => {
    const requestUser = new UserNameRequest();
    requestUser.setUserName(userName);
    const response = await this.props.apiClient.getUserWithUsername(
      requestUser
    );
    const user = response.toObject();
    const profile = { facesList: [], weight: 0, user: user };
    this.props.onProfileSelect(profile);
  };
  render() {
    return (
      <div className="container-popup">
        <div className="column-popup">
          <H2>We couldn't find you. Please enter your instagram username.</H2>
          <Form onSubmit={this.handleSubmit} />
          <br />
          <Button onClick={this.props.goToExample}>
            I don't use instagram
          </Button>
        </div>
      </div>
    );
  }
}

SearchProfile.propTypes = {
  history: PropTypes.shape({
    push: PropTypes.shape({
      pathname: PropTypes.string,
      state: PropTypes.object
    })
  })
};

export default withRouter(SearchProfile);
