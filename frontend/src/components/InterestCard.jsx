import React from "react";
import PropTypes from "prop-types";
import Slideshow from "./Slideshow";
import InterestFooter from "./InterestFooter";

const InterestCard = props => {
  return (
    <div className="dashboardCard">
      <div className="cardGrid">
        <Slideshow slides={props.slides} />
        <InterestFooter title={props.title} details={props.details} />
      </div>
    </div>
  );
};

InterestCard.propTypes = {
  slides: PropTypes.array,
  title: PropTypes.string,
  details: PropTypes.string
};

export default InterestCard;
