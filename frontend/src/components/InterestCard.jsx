import React, { Component } from "react";
import Slideshow from "./Slideshow";
import InterestFooter from "./InterestFooter";

class InterestCard extends Component {
  render() {
    const { slides, title, details } = this.props;

    return (
      <div className="dashboardCard">
        <div className="cardGrid">
          <Slideshow slides={slides} />
          <InterestFooter title={title} details={details} />
        </div>
      </div>
    );
  }
}

export default InterestCard;
