import React, { Component } from "react";
import InterestCard from "../components/InterestCard";
import "./../Dashboard.css";
import ProfileCard from "../components/ProfileCard";

class Dashboard extends Component {
  render() {
    return (
      <div className="dashboard">
        <ProfileCard pictureUrl="https://avatars3.githubusercontent.com/u/17454617?s=460&v=4" />
        <InterestCard />
      </div>
    );
  }
}

export default Dashboard;
