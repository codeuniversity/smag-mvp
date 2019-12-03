import React, { Component } from "react";

function Button(props) {
  if (props.buttonlink) {
    return (
      <div className="button">
        <a href={props.buttonlink}>{props.children}</a>
      </div>
    );
  }

  return (
    <div className="button" onClick={props.onClick}>
      {props.children}
    </div>
  );
}

export default Button;
