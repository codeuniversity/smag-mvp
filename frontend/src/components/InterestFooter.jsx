import React from "react";
import PropTypes from "prop-types";

const InterestFooter = props => {
  return (
    <div className="interestFooter">
      <h2>{props.title}</h2>
      <p>{props.details}</p>
    </div>
  );
};

InterestFooter.propTypes = {
  title: PropTypes.string
};

export default InterestFooter;
