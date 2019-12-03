import React from "react";
import PropTypes from "prop-types";

const getDetailsString = details => {
  let detailsString = "";

  details.map((item, index) => {
    detailsString += item;

    if (index < details.length - 1) {
      detailsString += ", ";
    }
  });

  return detailsString;
};

const InterestFooter = props => {
  return (
    <div className="interestFooter">
      <h2>{props.title}</h2>
      <p>{getDetailsString(props.details)}</p>
    </div>
  );
};

InterestFooter.propTypes = {
  title: PropTypes.string
};

export default InterestFooter;
