import React, { Component } from "react";

function H1(props) {
  return (
    <div className="headline">
      <h1>{props.children}</h1>
    </div>
  );
}

export default H1;
