import React, { Component, useState } from "react";
import { withRouter, history } from "react-router";
import Button from "../components/Button";
import H1 from "../components/H1";
import H2 from "../components/H2";
import {
  User,
  UserNameRequest,
  UserIdRequest,
  Post,
  UserIdResponse,
  UserSearchResponse
} from "../protofiles/usersearch_pb.js";
import { UserSearchServiceClient } from "../protofiles/usersearch_grpc_web_pb";
async function fetchExampleProfile(apiClient) {
  const requestUser = new UserNameRequest();

  requestUser.setUserName("codeuniversity");
  const response = await apiClient.getUserWithUsername(requestUser);
  const user = response.toObject();
  return user;
}
function ExampleProfileSelection(props) {
  return (
    <div className="container-popup">
      <div className="column-popup">
        <H1>We couldn't find any data of you.</H1>
        <br />
        <H2>
          Would you like to continue the exhibition with the profile of CODE?
        </H2>
        <div className="column-one-fourth" />
        <div className="column-one-fourth">
          <Button onClick={props.goToEnd}>End</Button>
        </div>
        <div className="column-one-fourth">
          <Button
            onClick={() =>
              fetchExampleProfile(props.apiClient).then(user => {
                const profile = { facesList: [], weight: 0, user: user };
                props.onProfileSelect(profile);
              })
            }
          >
            Yes
          </Button>
        </div>
        <div className="column-one-fourth" />
      </div>
    </div>
  );
}

export default ExampleProfileSelection;
