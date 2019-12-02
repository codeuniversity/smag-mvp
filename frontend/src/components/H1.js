import React, { Component } from "react";

class H1 extends Component {
  render() {
    return (
      <div className="headline">
        <h1>{this.props.h1}</h1>
      </div>
    );
  }
}

export default H1;
