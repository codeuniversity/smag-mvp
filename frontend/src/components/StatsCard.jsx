import React from "react";
import PropTypes from "prop-types";

const StatsCard = props => {
  return (
    <div className="dashboardCard statsCard dashboard-stats">
      <p>We were able to reuse</p>
      <h1>{props.count}</h1>
      <p>snippets of your data.</p>
    </div>
  );
};

StatsCard.propTypes = {
  count: PropTypes.number
};

export default StatsCard;
