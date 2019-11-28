import React from "react";
import { CameraFeed } from "./camera-feed";
import Title from "./Title";
import { FaceSearchRequest } from "../protofiles/usersearch_pb";
import { UserSearchServicePromiseClient } from "../protofiles/usersearch_grpc_web_pb";

const uploadImage = async file => {
  var reader = new FileReader();
  reader.onloadend = async () => {
    const dataUrl = reader.result;
    const base64Data = dataUrl.split(",")[1];
    console.log(base64Data);
    const client = new UserSearchServicePromiseClient("http://localhost:4000");
    const request = new FaceSearchRequest();
    request.setBase64encodedpicture(base64Data);
    const response = await client.searchSimilarFaces(request);
    console.log(response);
    debugger;
  };
  reader.readAsDataURL(file);
};

function Start(props) {
  return (
    <div className="body">
      <div className="container">
        <div className="column-center">
          <Title />
          <p> Take a pictures! </p>
          <CameraFeed sendFile={uploadImage} />
        </div>
      </div>
    </div>
  );
}

export default Start;
