import React, { Component } from "react";
import H2 from "./H2";

class LocationCard extends Component {
  render() {
    return (
      <div className="smallCard bioCard dashboard-location">
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

export default LocationCard;
