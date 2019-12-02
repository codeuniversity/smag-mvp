import React, { useState, useEffect } from "react";
import Start from "./Start";
import ProfileSelection from "./ProfileSelection";
import { UserSearchServicePromiseClient } from "../protofiles/usersearch_grpc_web_pb";
import "../creativeCode.css";

const START_PAGE = "start";
const PROFILE_SELECTION_PAGE = "profile_selection";
const NECESARY_FACE_SAMPLES = 5;
const apiClient = new UserSearchServicePromiseClient("http://localhost:4000");

function FlowStateWrapper(props) {
  const [page, setPage] = useState(START_PAGE);
  const { faceHits, addFaceHits } = useFaceHitState(() =>
    setPage(PROFILE_SELECTION_PAGE)
  );

  switch (page) {
    case START_PAGE:
      return (
        <Start
          apiClient={apiClient}
          faceHits={faceHits}
          addFaceHits={addFaceHits}
        />
      );
    case PROFILE_SELECTION_PAGE:
      return <ProfileSelection apiClient={apiClient} faceHits={faceHits} />;
    default:
      return `page ${page} not found`;
  }
}

function FlowWrapper(props) {
  return (
    <>
      <div className="white-background"></div>
      <div className="container">
        <FlowStateWrapper />
      </div>
    </>
  );
}

export default FlowWrapper;

function useFaceHitState(onEnoughFacesCollected) {
  const [faceHits, setFaceHits] = useState({});
  const [faceSampleCount, setFaceSampleCount] = useState(0);
  const addFaceHits = faces => {
    setFaceSampleCount(prevCount => prevCount + 1);
    setFaceHits(prevHits => mergeFacesIntoHits(faces, prevHits));
  };

  useEffect(() => {
    if (faceSampleCount >= NECESARY_FACE_SAMPLES) {
      onEnoughFacesCollected();
    }
  }, [faceSampleCount]);

  return {
    faceHits,
    addFaceHits
  };
}

function mergeFacesIntoHits(faces, hits) {
  const newHits = { ...hits };
  faces.forEach(face => {
    newHits[face.postId] = [...(newHits[face.postId] || []), face];
  });

  return newHits;
}