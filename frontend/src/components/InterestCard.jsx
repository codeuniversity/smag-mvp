import React, { Component } from "react";
import Slideshow from "./Slideshow";
import InterestFooter from "./InterestFooter";

class InterestCard extends Component {
  constructor(props) {
    super(props);

    this.state = {
      slides: this.props.slides,
      title: this.props.title,
      details: this.props.details
    };
  }

  render() {
    return (
      <div className="dashboardCard">
        <div className="cardGrid">
          <Slideshow slides={this.state.slides} />
          <InterestFooter
            title={this.state.title}
            details={this.state.details}
          />
        </div>
      </div>
    );
  }
}

export default InterestCard;
