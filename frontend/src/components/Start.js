import React, { useState } from "react";
import { CameraFeed } from "./camera-feed";
import Title from "./Title";
import { FaceSearchRequest } from "../protofiles/usersearch_pb";
import { UserSearchServicePromiseClient } from "../protofiles/usersearch_grpc_web_pb";
import IGPost from "./IGPost";
import "../creativeCode.css"
function Start(props) {
  const [similarFaces, setSimilarFaces] = useState([]);

  const uploadImage = async file => {
    const reader = new FileReader();
    reader.onloadend = async () => {
      const dataUrl = reader.result;
      const base64Data = dataUrl.split(",")[1];
      console.log(base64Data);
      const client = new UserSearchServicePromiseClient(
        "http://localhost:4000"
      );
      const request = new FaceSearchRequest();
      request.setBase64encodedpicture(base64Data);
      const response = await client.searchSimilarFaces(request);
      console.log(response);
      const faces = response.getFacesList().map(protoFace => ({
        postId: protoFace.getPostId(),
        x: protoFace.getX(),
        y: protoFace.getY(),
        width: protoFace.getWidth(),
        height: protoFace.getHeight(),
        fullImageSrc: protoFace.getFullImageSrc()
      }));

      setSimilarFaces(faces);
    };
    reader.readAsDataURL(file);
  };

  return (
    <div className="body">
      <div className="white-background"></div>
      <div className="container">
        <div className="column-center">
          <Title />
          <p> Take a pictures! </p>
          <CameraFeed sendFile={uploadImage} />

          {similarFaces.map(face => (
            <IGPost post={{ img: face.fullImageSrc, shortcode: "" }} />
          ))}
        </div>
      </div>
    </div>
  );
}

export default Start;
