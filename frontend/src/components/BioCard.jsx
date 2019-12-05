import React, { Component } from "react";
import H2 from "./H2";

class BioCard extends Component {
  render() {
    return (
      <div className="dashboardCard bioCard">
        <table>
          <tr>
            <th>
              <h2 className="dashboard-headline">{this.props.bio}</h2>
            </th>
          </tr>
        </table>
      </div>
    );
  }
}

export default BioCard;
