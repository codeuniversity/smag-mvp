import React from "react";
import PropTypes from "prop-types";

const ProfileCard = props => {
  return (
    <div className="dashboardCard">
      <img
        src={props.pictureUrl}
        alt={props.alt}
        className="profilePictureImg"
      />
    </div>
  );
};

ProfileCard.propTypes = {
  pictureUrl: PropTypes.string
};

export default ProfileCard;
