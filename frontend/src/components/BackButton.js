import React, { Component, useState } from "react";
import { withRouter, history } from "react-router";
import backicon from "./img/chevron-left-solid.svg";

const BackButton = () => (
  <div className="BackButton">
    <a href="/">
      <img className="BackButtonImg" src={backicon} width="24px" />
    </a>
  </div>
);

export default BackButton;
