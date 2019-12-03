import React, { useEffect, useState } from "react";
import {
  WeightedPosts,
  WeightedPostWithFaces,
  Face
} from "../protofiles/usersearch_pb";
import IGPost from "./IGPost";
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

function ProfileSelection({ apiClient, faceHits, nextPage }) {
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
  console.log(foundProfiles);
  return (
    <div className="body">
      <div className="column-center">
        {foundProfiles.map(profile => (
          <div className="container-card" key={profile.user.userName}>
            <p>
              {profile.user.userName}:{" "}
              {Math.round((profile.weight / weightSum) * 100)}% confidence
              <button onClick={nextPage}>Next Page</button>
            </p>
            {uniqWith(profile.facesList, (a, b) => a.postId === b.postId).map(
              face => (
                <IGPost
                  key={`${profile.user.userName}/${face.postId}`}
                  post={{ img: face.fullImageSrc, shortcode: "" }}
                />
              )
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

export default ProfileSelection;
