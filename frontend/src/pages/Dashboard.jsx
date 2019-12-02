import React, { Component } from "react";
import InterestCard from "../components/InterestCard";
import "./../Dashboard.css";
import ProfileCard from "../components/ProfileCard";
import StatsCard from "../components/StatsCard";
import BioCard from "../components/BioCard";

class Dashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      slides0: [
        "https://www.diabetes.org/sites/default/files/styles/full_width/public/2019-06/Healthy%20Food%20Made%20Easy%20-min.jpg",
        "http://www.islandsinthesun.com/img/home-maldives5.jpg",
        "https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/best-running-shoes-lead-02-1567016766.jpg?crop=0.502xw:1.00xh;0.0577xw,0&resize=640:*"
      ],
      slides1: [
        "https://media.pitchfork.com/photos/5da53402163e3300096f6a07/2:1/w_790/Coldplay.jpg",
        "https://townsquare.media/site/366/files/2017/06/Rammstein_2017-1-7.jpg?w=980&q=75"
      ],
      slides2: [
        "https://img.chefkoch-cdn.de/rezepte/27131006360850/bilder/1180995/crop-600x400/meine-koenigsberger-klopse.jpg",
        "https://www.omoxx.com/wp-content/uploads/2018/05/zucchini-auberginen-pasta.jpg",
        "https://media-cdn.tripadvisor.com/media/photo-s/10/78/00/4c/pizza-peperoni-wurst.jpg"
      ],
      slides3: [
        "https://images.unsplash.com/photo-1530938959149-d8eb57633f2c?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2850&q=80"
      ],
      slides4: [
        "https://d13ezvd6yrslxm.cloudfront.net/wp/wp-content/images/sister-act-3-update.jpg"
      ]
    };
  }

  render() {
    return (
      <div className="dashboard">
        <h1 className="dashboardTitle">Here's what we found out about you:</h1>
        <div className="dashboardGrid">
          <ProfileCard pictureUrl="https://avatars3.githubusercontent.com/u/17454617?s=460&v=4" />
          <InterestCard
            title="Favorites"
            details={["Holidays", "Shoes", "Food"]}
            slides={this.state.slides0}
          />
          <InterestCard
            title="Music"
            details={["Coldplay", "Nickelback", "Rammstein"]}
            slides={this.state.slides1}
          />
          <StatsCard count="53" />
          <BioCard />
          <InterestCard
            title="Food"
            details={["Pizza", "Pasta", "Königsberger Klopse"]}
            slides={this.state.slides2}
          />
          <InterestCard
            title="Brands"
            details={["Nike"]}
            slides={this.state.slides3}
          />
          <InterestCard
            title="Friends & Family"
            details={["Sister, Mother"]}
            slides={this.state.slides4}
          />
        </div>
      </div>
    );
  }
}

export default Dashboard;
