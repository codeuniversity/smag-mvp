import React, { useState, useEffect } from "react";
import Start from "./Start";
import ProfileSelection from "./ProfileSelection";
import { UserSearchServicePromiseClient } from "../protofiles/usersearch_grpc_web_pb";
import Dashboard from "../pages/Dashboard";
import Network from "../pages/Network";
import "../creativeCode.css";
import Greeting from "../pages/Greeting";
import ExampleProfileSelection from "../pages/ExampleProfileSelection";
import SearchProfile from "../pages/SearchProfile";
import GroupIntent from "../pages/GroupIntent";
import EndScreen from "../pages/endscreen";

const GREETING_PAGE = "greeting";
const START_PAGE = "start";
const PROFILE_SELECTION_PAGE = "profile_selection";
const DASHBOARD_PAGE = "dashboard";
const EXAMPLE_PROFILE_PAGE = "example-profile";
const SEARCH_PROFILE_PAGE = "search-profile";
const GROUP_INTENT_PAGE = "group_intent";
const NETWORK_PAGE = "network";
const END_PAGE = "endscreen";
const NECESARY_FACE_SAMPLES = 5;
const apiClient = new UserSearchServicePromiseClient("http://localhost:4000");

function FlowStateWrapper(props) {
  const [page, setPage] = useState(GREETING_PAGE);
  const { faceHits, addFaceHits, faceSampleCount } = useFaceHitState(() =>
    setPage(PROFILE_SELECTION_PAGE)
  );
  const [selectedProfile, setSelectedProfile] = useState(null);

  switch (page) {
    case GREETING_PAGE:
      return <Greeting nextPage={() => setPage(START_PAGE)} />;
    case START_PAGE:
      return (
        <Start
          apiClient={apiClient}
          faceHits={faceHits}
          addFaceHits={addFaceHits}
          progress={faceSampleCount / NECESARY_FACE_SAMPLES}
        />
      );
    case PROFILE_SELECTION_PAGE:
      return (
        <ProfileSelection
          apiClient={apiClient}
          faceHits={faceHits}
          goToSearch={() => setPage(SEARCH_PROFILE_PAGE)}
          goToExample={() => setPage(EXAMPLE_PROFILE_PAGE)}
          onProfileSelect={profile => {
            setPage(DASHBOARD_PAGE);
            setSelectedProfile(profile);
          }}
        />
      );
    case SEARCH_PROFILE_PAGE:
      return (
        <SearchProfile
          apiClient={apiClient}
          goToExample={() => setPage(EXAMPLE_PROFILE_PAGE)}
          onProfileSelect={profile => {
            setSelectedProfile(profile);
            setPage(DASHBOARD_PAGE);
          }}
        />
      );
    case EXAMPLE_PROFILE_PAGE:
      return (
        <ExampleProfileSelection
          apiClient={apiClient}
          goToEnd={() => setPage(END_PAGE)}
          onProfileSelect={profile => {
            setSelectedProfile(profile);
            setPage(DASHBOARD_PAGE);
          }}
        />
      );
    case DASHBOARD_PAGE:
      return (
        <Dashboard
          profile={selectedProfile}
          apiClient={apiClient}
          nextPage={() => setPage(NETWORK_PAGE)}
        />
      );
    case NETWORK_PAGE:
      return (
        <Network
          profile={selectedProfile}
          apiClient={apiClient}
          nextPage={() => setPage(GROUP_INTENT_PAGE)}
        />
      );
    case GROUP_INTENT_PAGE:
      return <GroupIntent nextPage={() => setPage(END_PAGE)} />;
    case END_PAGE:
      return <EndScreen nextPage={() => window.location.reload()} />;
    default:
      return `page ${page} not found`;
  }
}

function FlowWrapper(props) {
  return (
    <>
      <FlowStateWrapper />
    </>
  );
}

export default FlowWrapper;

function useFaceHitState(onEnoughFacesCollected) {
  const [faceHits, setFaceHits] = useState({});
  const [faceSampleCount, setFaceSampleCount] = useState(0);
  const [hasTriggerEnoughCollected, setHasTriggerEnoughCollected] = useState(
    false
  );
  const addFaceHits = faces => {
    if (faces.length == 0) return;

    setFaceSampleCount(prevCount => prevCount + 1);
    setFaceHits(prevHits => mergeFacesIntoHits(faces, prevHits));
  };

  useEffect(() => {
    if (
      faceSampleCount >= NECESARY_FACE_SAMPLES &&
      !hasTriggerEnoughCollected
    ) {
      setHasTriggerEnoughCollected(true);
      onEnoughFacesCollected();
    }
  }, [faceSampleCount]);

  return {
    faceHits,
    addFaceHits,
    faceSampleCount
  };
}

function mergeFacesIntoHits(faces, hits) {
  const newHits = { ...hits };
  faces.forEach(face => {
    newHits[face.postId] = [...(newHits[face.postId] || []), face];
  });

  return newHits;
}
