import React, { Component } from "react";
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
            <H1>Welcome to SocialRecord</H1>
            <H2>Sit back and enjoy the experience.</H2>
            <div style={{ marginTop: 50 }}>
              <Button onClick={this.props.nextPage}>Start</Button>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default Greeting;
