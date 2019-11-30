import React, { Component } from "react";
import InterestCard from "../components/InterestCard";
import "./../Dashboard.css";

class Dashboard extends Component {
  render() {
    return (
      <div className="dashboard">
        <InterestCard />
      </div>
    );
  }
}

export default Dashboard;
