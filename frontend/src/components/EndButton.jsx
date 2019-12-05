import React, { Component } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTimes } from "@fortawesome/free-solid-svg-icons";
import { Link } from "react-router-dom";
import "./../css/endButton.css";

class EndButton extends Component {
  render() {
    return (
      <div className="endButton">
        <Link to={this.props.link}>
          <FontAwesomeIcon icon={faTimes} />
        </Link>
      </div>
    );
  }
}

export default EndButton;
