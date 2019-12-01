import React, { Component } from "react";

class InterestFooter extends Component {
  getDetailsString = () => {
    let detailsString = "";

    this.props.details.map((item, index) => {
      detailsString += item;

      if (index < this.props.details.length - 1) {
        detailsString += ", ";
      }
    });

    return detailsString;
  };

  render() {
    return (
      <div className="interestFooter">
        <h2>{this.props.title}</h2>
        <p>{this.getDetailsString()}</p>
      </div>
    );
  }
}

export default InterestFooter;
