import React, { Component } from "react";

class Button extends Component {
  render() {
    return (
      <div className="button">
        <a href={this.props.buttonlink}>{this.props.button}</a>
      </div>
    );
  }
}

export default Button;
