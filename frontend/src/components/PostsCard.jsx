import React from "react";
import PropTypes from "prop-types";
import Slideshow from "./Slideshow";

const PostsCard = props => {
  return (
    <div className="dashboardCard">
      <Slideshow slides={props.slides} />
    </div>
  );
};

PostsCard.propTypes = {
  slides: PropTypes.array
};

export default PostsCard;
