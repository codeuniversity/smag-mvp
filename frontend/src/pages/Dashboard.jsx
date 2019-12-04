import React, { useState, useEffect } from "react";
import InterestCard from "../components/InterestCard";
import "./../Dashboard.css";
import ProfileCard from "../components/ProfileCard";
import StatsCard from "../components/StatsCard";
import BioCard from "../components/BioCard";
import { UserIdRequest } from "../protofiles/usersearch_grpc_web_pb";
import Button from "../components/Button";

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

function Dashboard({ profile, apiClient }) {
  const [posts, setPosts] = useState([]);
  const [dataPointCount, setDataPointCount] = useState(null);

  useEffect(() => {
    fetchPosts(apiClient, profile.user.id).then(setPosts);
    fetchDataPoints(apiClient, profile.user.id).then(setDataPointCount);
  }, []);

  const slides0 = profile.facesList.map(face => face.fullImageSrc);
  const slides1 = posts.map(post => post.imgUrl);

  const slides2 = [
    "https://img.chefkoch-cdn.de/rezepte/27131006360850/bilder/1180995/crop-600x400/meine-koenigsberger-klopse.jpg",
    "https://www.omoxx.com/wp-content/uploads/2018/05/zucchini-auberginen-pasta.jpg",
    "https://media-cdn.tripadvisor.com/media/photo-s/10/78/00/4c/pizza-peperoni-wurst.jpg"
  ];
  const slides3 = [
    "https://images.unsplash.com/photo-1530938959149-d8eb57633f2c?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2850&q=80"
  ];
  const slides4 = [
    "https://d13ezvd6yrslxm.cloudfront.net/wp/wp-content/images/sister-act-3-update.jpg"
  ];

  return (
    <div className="dashboard">
      <h1 className="dashboardTitle">Here's what we found out about you:</h1>
      <div className="dashboardGrid">
        <ProfileCard
          pictureUrl={profile.user.avatarUrl}
          alt={profile.facesList[0] && profile.facesList[0].fullImageSrc}
        />
        <InterestCard
          title="Favorites"
          details={["Holidays", "Shoes", "Food"]}
          slides={slides0}
        />
        <InterestCard
          title="Music"
          details={["Coldplay", "Nickelback", "Rammstein"]}
          slides={slides1}
        />
        <StatsCard count={dataPointCount} />
        <BioCard bio={profile.user.bio} />
        <InterestCard
          title="Food"
          details={["Pizza", "Pasta", "Königsberger Klopse"]}
          slides={slides2}
        />
        <InterestCard title="Brands" details={["Nike"]} slides={slides3} />
        <InterestCard
          title="Friends & Family"
          details={["Sister, Mother"]}
          slides={slides4}
        />
      </div>
      <div className="dashboardFooter">
        <Button buttonlink="/endscreen">
          See more details about your network.
        </Button>
      </div>
    </div>
  );
}

export default Dashboard;
