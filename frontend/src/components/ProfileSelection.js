import React, { useEffect, useState } from "react";
import {
  WeightedPosts,
  WeightedPostWithFaces,
  Face
} from "../protofiles/usersearch_pb";
import IGPost from "./IGPost";
import H1 from "./H1";
import Button from "./Button";
import uniqWith from "lodash/uniqWith";
import FaceHitAnimation from "./FaceHitAnimation";

async function searchProfiles(apiClient, faceHits) {
  const weightedPosts = new WeightedPosts();
  const posts = Object.entries(faceHits).map(([postId, faces]) => {
    const postWithFaces = new WeightedPostWithFaces();
    postWithFaces.setPostId(parseInt(postId));
    postWithFaces.setWeight(faces.length);

    postWithFaces.setFacesList(
      faces.map(face => {
        const protoFace = new Face();
        protoFace.setX(face.x);
        protoFace.setY(face.y);
        protoFace.setWidth(face.width);
        protoFace.setHeight(face.height);
        protoFace.setPostId(face.postId);
        protoFace.setFullImageSrc(face.fullImageSrc);

        return protoFace;
      })
    );

    return postWithFaces;
  });
  weightedPosts.setPostsList(posts);
  const weightedUsers = await apiClient.searchUsersWithWeightedPosts(
    weightedPosts
  );
  const profiles = weightedUsers
    .getUsersWithFacesList()
    .map(userWithFaces => userWithFaces.toObject())
    .sort((a, b) => b.weight - a.weight);

  return profiles;
}

function ProfileSelection({
  apiClient,
  faceHits,
  onProfileSelect,
  goToExample,
  goToSearch
}) {
  const [foundProfiles, setFoundProfiles] = useState([]);
  const [loadingAnimationDone, setLoadingAnimationDone] = useState(false);
  useEffect(() => {
    searchProfiles(apiClient, faceHits).then(profiles => {
      setFoundProfiles(profiles);
    });
  }, []);

  if (foundProfiles.length == 0 || !loadingAnimationDone) {
    return (
      <FaceHitAnimation
        faceHits={faceHits}
        onAnimationFinished={() => setLoadingAnimationDone(true)}
      />
    );
  }

  const weightSum = foundProfiles.reduce(
    (sum, profile) => sum + profile.weight,
    0
  );

  return (
    <div className="container-popup">
      <div className="column-popup">
        <H1>Please select your profile</H1>

        {foundProfiles.slice(0, 6).map(profile => (
          <div
            className="profile-card"
            key={profile.user.userName}
            onClick={() => onProfileSelect(profile)}
          >
            <div className="avatar-image-container">
              <img className="avatar-image" src={profile.user.avatarUrl} />
            </div>

            <div
              className="sub-headline-profile"
              style={{
                paddingLeft: 10
              }}
            >
              <div>
                {profile.user.userName}: <br />
              </div>
              <div style={{ fontSize: 14 }}>
                {Math.round((profile.weight / weightSum) * 100)}% confidence
              </div>
            </div>
          </div>
        ))}
        <div className="profile-button">
          <div className="container-profile">
            <div className="column-one-fourth">
              <Button onClick={goToSearch}>My profile is not shown.</Button>
            </div>
            <div className="column-one-fourth">
              <Button onClick={goToExample} style={{ fontSize: 20 }}>
                I don't use instagram.
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default ProfileSelection;
