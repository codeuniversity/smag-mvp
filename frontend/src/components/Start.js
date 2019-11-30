import React, { useState } from "react";
import { CameraFeed } from "./camera-feed";
import Title from "./Title";
import { FaceSearchRequest } from "../protofiles/usersearch_pb";
import { UserSearchServicePromiseClient } from "../protofiles/usersearch_grpc_web_pb";
import IGPost from "./IGPost";
import "../creativeCode.css";

function findFacesInImage(onFindFaces) {
  return async file => {
    const reader = new FileReader();
    reader.onloadend = async () => {
      const dataUrl = reader.result;
      const base64Data = dataUrl.split(",")[1];
      const client = new UserSearchServicePromiseClient(
        "http://localhost:4000"
      );
      const request = new FaceSearchRequest();
      request.setBase64encodedpicture(base64Data);
      const response = await client.searchSimilarFaces(request);
      const faces = response.getFacesList().map(protoFace => ({
        postId: protoFace.getPostId(),
        x: protoFace.getX(),
        y: protoFace.getY(),
        width: protoFace.getWidth(),
        height: protoFace.getHeight(),
        fullImageSrc: protoFace.getFullImageSrc()
      }));

      onFindFaces(faces);
    };
    reader.readAsDataURL(file);
  };
}

function mergeFacesIntoHits(faces, hits) {
  const newHits = { ...hits };
  faces.forEach(face => {
    newHits[face.postId] = [...(newHits[face.postId] || []), face];
  });

  return newHits;
}

function Start(props) {
  const [similarFaces, setSimilarFaces] = useState({});

  const onFileSubmit = findFacesInImage(faces =>
    setSimilarFaces(prevHits => mergeFacesIntoHits(faces, prevHits))
  );
  return (
    <div className="body">
      <div className="white-background"></div>
      <div className="container">
        <div className="column-center">
          <Title />
          <p>Take a pictures!</p>
          <CameraFeed onFileSubmit={onFileSubmit} />

          {Object.entries(similarFaces)
            .sort((a, b) => b[1].length - a[1].length)
            .map(([postId, faces]) => {
              console.log(faces.length);
              return (
                <IGPost
                  key={postId}
                  post={{ img: faces[0].fullImageSrc, shortcode: "" }}
                />
              );
            })}
        </div>
      </div>
    </div>
  );
}

export default Start;
