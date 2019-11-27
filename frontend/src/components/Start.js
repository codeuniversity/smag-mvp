import React, { Fragment } from "react";
import ReactDOM from "react-dom";
// import axios from "axios";
import { withRouter } from "react-router";
import { CameraFeed } from "./Camera-feed";
import Title from "./Title";

// Upload to local seaweedFS instance
const uploadImage = async file => {
  const formData = new FormData();
  formData.append("file", file);

  // Connect to a seaweedfs instance
};

function Start() {
  return (
    <div className="body">
      <div className="container">
        <div className="column-center">
          <Title />
          <p>Take a pictures!</p>
          <CameraFeed onFileSubmit={uploadImage} />
        </div>
      </div>
    </div>
  );
}

export default withRouter(Start);
