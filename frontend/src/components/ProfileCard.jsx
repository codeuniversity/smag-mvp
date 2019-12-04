import React from "react";
import PropTypes from "prop-types";

const ProfileCard = props => {
  return (
    <div className="dashboardCard">
      <img
        src={this.props.pictureUrl}
        alt={this.props.alt}
        className="profilePictureImg"
      />
    </div>
  );
};

ProfileCard.propTypes = {
  pictureUrl: PropTypes.string
};

export default ProfileCard;
