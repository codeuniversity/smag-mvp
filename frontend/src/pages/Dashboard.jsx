import React, { useState, useEffect } from "react";
import InterestCard from "../components/InterestCard";
import "./../Dashboard.css";
import ProfileCard from "../components/ProfileCard";
import StatsCard from "../components/StatsCard";
import BioCard from "../components/BioCard";
import { UserIdRequest } from "../protofiles/usersearch_grpc_web_pb";
import Button from "../components/Button";
import EndButton from "../components/EndButton";
import uniqWith from "lodash/uniqWith";
import PostsCard from "../components/PostsCard";

async function fetchPosts(apiClient, userId) {
  const userIdRequest = new UserIdRequest();
  userIdRequest.setUserId(userId);
  const response = await apiClient.getInstaPostsWithUserId(userIdRequest);
  const posts = response.getInstaPostsList();

  return posts.map(post => post.toObject()).filter(post => !!post.imgUrl);
}

async function fetchDataPoints(apiClient, userId) {
  const userIdRequest = new UserIdRequest();
  userIdRequest.setUserId(userId);

  const response = await apiClient.dataPointCountForUserId(userIdRequest);
  return response.getCount();
}

async function fetchCities(apiClient, userId) {
  const userIdRequest = new UserIdRequest();
  userIdRequest.setUserId(userId);

  const response = await apiClient.findCitiesForUserId(userIdRequest);
  return response.getCityNamesList();
}

function Dashboard({ profile, apiClient, nextPage }) {
  const [posts, setPosts] = useState([]);
  const [dataPointCount, setDataPointCount] = useState(null);
  const [foundCities, setFoundCities] = useState([]);
  useEffect(() => {
    fetchPosts(apiClient, profile.user.id).then(setPosts);
    fetchDataPoints(apiClient, profile.user.id).then(setDataPointCount);
    fetchCities(apiClient, profile.user.id).then(setFoundCities);
  }, []);

  // get the imageSrc of the de-duplicated list of posts where we found that face
  const slides0 = uniqWith(
    profile.facesList,
    (a, b) => a.postId === b.postId
  ).map(face => face.fullImageSrc);
  const postImageArray = posts.map(post => post.imgUrl);

  return (
    <div className="dashboard">
      <EndButton link="/" />
      <h1 className="dashboardTitle">Here's what we found out about you:</h1>
      <div className="dashboardGrid">
        <ProfileCard
          className="dashboard-avatar"
          pictureUrl={profile.user.avatarUrl}
          alt={profile.facesList[0] && profile.facesList[0].fullImageSrc}
        />
        <PostsCard
          className="dashboard-posts"
          title="Your images"
          description=""
          slides={postImageArray}
        />
        <StatsCard className="dashboard-stats" count={dataPointCount} />
        <BioCard className="dashboard-bio" bio={profile.user.bio} />
        <InterestCard
          title="Cities you were in"
          details={
            foundCities.length > 0
              ? foundCities.join(", ")
              : "We didn't find any cities for you."
          }
          slides={[
            "http://www.brinknews.com/wp-content/uploads/2019/09/GettyImages-854625740.jpg"
          ]}
        />
      </div>
      <div className="dashboardFooter">
        <Button onClick={nextPage}>
          More details about your personal network
        </Button>
      </div>
    </div>
  );
}

export default Dashboard;
