import React, { useState } from "react";
import { CameraFeed } from "./camera-feed";
import H1 from "./H1";
import { FaceSearchRequest } from "../protofiles/usersearch_pb";
import IGPost from "./IGPost";

function findFacesInImage(apiClient, onFindFaces) {
  return async file => {
    const reader = new FileReader();
    reader.onloadend = async () => {
      const dataUrl = reader.result;
      const base64Data = dataUrl.split(",")[1];

      const request = new FaceSearchRequest();
      request.setBase64encodedpicture(base64Data);
      const response = await apiClient.searchSimilarFaces(request);
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

function Start({ apiClient, faceHits, addFaceHits, nextPage }) {
  const onFileSubmit = findFacesInImage(apiClient, addFaceHits);

  return (
    <div className="body">
      <div className="column-center">
        <H1 />
        <p>Take a pictures!</p>
        <CameraFeed onFileSubmit={onFileSubmit} />

        {Object.entries(faceHits)
          .sort((a, b) => b[1].length - a[1].length)
          .map(([postId, faces]) => {
            return (
              <IGPost
                key={postId}
                post={{ img: faces[0].fullImageSrc, shortcode: "" }}
              />
            );
          })}
      </div>
    </div>
  );
}

export default Start;
