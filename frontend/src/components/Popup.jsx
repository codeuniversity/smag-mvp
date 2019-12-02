import React, { Component, useState } from "react";
import { withRouter, history } from "react-router";
import Button from "../components/Button";
import H1 from "../components/H1";
import "./../index.css";

function Popup({ post }) {
  return (
    <div className="container-popup">
      <div className="column-popup">
        <H1 h1="Is this your profile?" />
        <img className="avatar-image" src="" />
        <div className="sub-headline">username</div>
        <div className="column-one-fourth" />
        <div className="column-one-fourth">
          <Button button="no" buttonlink="/" />
        </div>
        <div className="column-one-fourth">
          <Button button="yes" buttonlink="/result" />
        </div>
        <div className="column-one-fourth" />
      </div>
    </div>
  );
}

export default Popup;
