import React, { Component } from "react";

class StatsCard extends Component {
  render() {
    return (
      <div className="dashboardCard statsCard">
        <p>We were able to reuse</p>
        <h1>{this.props.count}</h1>
        <p>snippets of your data.</p>
      </div>
    );
  }
}

export default StatsCard;
