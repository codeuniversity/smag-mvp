import React, { Component, useState } from "react";
import { withRouter, history } from "react-router";
import Button from "../components/Button";
import H1 from "../components/H1";

function Popup() {
  return (
    <div className="container-popup">
      <div className="column-popup">
        <H1>Is this your profile?</H1>
        <img className="avatar-image" src="" />
        <div className="sub-headline">username</div>
        <div className="column-one-fourth" />
        <div className="column-one-fourth">
          <Button buttonlink="/">No</Button>
        </div>
        <div className="column-one-fourth">
          <Button buttonlink="/result">Yes</Button>
        </div>
        <div className="column-one-fourth" />
      </div>
    </div>
  );
}

export default Popup;
