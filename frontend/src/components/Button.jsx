import React, { Component } from "react";

function Button(props) {
  return (
    <div className="button">
      <a href={props.buttonlink}>{props.children}</a>
    </div>
  );
}

export default Button;
