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
            <H1>
              Think again about whether all this data about you should be
              visible to everyone!
              <br />
              If not, we want to give you 3 important tips.
            </H1>
            <H2>
              1.Think twice about what information you want to make public.
            </H2>
            <H2>2. Check your private settings again.</H2>
            <H2>3. Switch your profile to private.</H2>
            <Button buttonlink="/">End ssession</Button>
          </div>
        </div>
      </div>
    );
  }
}

export default Greeting;
