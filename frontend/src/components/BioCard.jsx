import React, { Component } from "react";
import H2 from "./H2";

class BioCard extends Component {
  render() {
    return (
      <div className="dashboardCard bioCard">
        <table>
          <tr>
            <th>
              <H2>{this.props.bio}</H2>
            </th>
          </tr>
        </table>
      </div>
    );
  }
}

export default BioCard;
