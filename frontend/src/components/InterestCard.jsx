import React, { Component } from "react";
import Slideshow from "./Slideshow";
import InterestFooter from "./InterestFooter";

class InterestCard extends Component {
  constructor(props) {
    super(props);

    this.state = {
      slides: [
        "https://www.diabetes.org/sites/default/files/styles/full_width/public/2019-06/Healthy%20Food%20Made%20Easy%20-min.jpg",
        "http://www.islandsinthesun.com/img/home-maldives5.jpg",
        "https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/best-running-shoes-lead-02-1567016766.jpg?crop=0.502xw:1.00xh;0.0577xw,0&resize=640:*"
      ]
    };
  }

  render() {
    return (
      <div className="dashboardCard">
        <div className="cardGrid">
          <Slideshow slides={this.state.slides} />
          <InterestFooter />
        </div>
      </div>
    );
  }
}

export default InterestCard;
