import React, { Component } from "react";
import { withRouter } from "react-router";
import Button from "../components/Button";
import "./../index.css";
import H1 from "../components/H1";
import H2 from "../components/H2";

class Greeting extends Component {
  render() {
    return (
      <div className="container-center">
        <div className="column-center">
          <div className="greeting">
            <H1 h1="Welcome to SocialRecord" />
            <H2 h2="Sit back and enjoy the experience." />
            <Button button="Start" buttonlink="/start" />
          </div>
        </div>
      </div>
    );
  }
}

export default Greeting;
