import React, { Component } from "react";

class ProfileCard extends Component {
  render() {
    return (
      <div className="dashboardCard">
        <img
          src={this.props.pictureUrl}
          alt="profile picture"
          className="profilePictureImg"
        />
      </div>
    );
  }
}

export default ProfileCard;
