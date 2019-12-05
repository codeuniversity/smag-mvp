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

function Dashboard({ profile, apiClient, nextPage }) {
  const [posts, setPosts] = useState([]);
  const [dataPointCount, setDataPointCount] = useState(null);

  useEffect(() => {
    fetchPosts(apiClient, profile.user.id).then(setPosts);
    fetchDataPoints(apiClient, profile.user.id).then(setDataPointCount);
  }, []);

  // get the imageSrc of the de-duplicated list of posts where we found that face
  const slides0 = uniqWith(
    profile.facesList,
    (a, b) => a.postId === b.postId
  ).map(face => face.fullImageSrc);
  const slides1 = posts.map(post => post.imgUrl);

  const food = [
    "https://img.chefkoch-cdn.de/rezepte/27131006360850/bilder/1180995/crop-600x400/meine-koenigsberger-klopse.jpg",
    "https://www.omoxx.com/wp-content/uploads/2018/05/zucchini-auberginen-pasta.jpg",
    "https://media-cdn.tripadvisor.com/media/photo-s/10/78/00/4c/pizza-peperoni-wurst.jpg"
  ];
  const favorites = [
    "http://www.islandsinthesun.com/img/home-maldives5.jpg",
    "https://www.diabetes.org/sites/default/files/styles/full_width/public/2019-06/Healthy%20Food%20Made%20Easy%20-min.jpg",
    "https://hips.hearstapps.com/hmg-prod.s3.amazonaws.com/images/best-running-shoes-lead-02-1567016766.jpg?crop=0.502xw:1.00xh;0.0577xw,0&resize=640:*"
  ];

  return (
    <div className="dashboard">
      <EndButton link="/" />
      <h1 className="dashboardTitle">Here's what we found out about you:</h1>
      <div className="dashboardGrid">
        <ProfileCard
          pictureUrl={profile.user.avatarUrl}
          alt={profile.facesList[0] && profile.facesList[0].fullImageSrc}
        />
        <InterestCard
          title="Your images"
          details={"These are the images where you found your face on."}
          slides={slides0}
        />
        <PostsCard slides={slides1} />
        <StatsCard count={dataPointCount} />
        <BioCard bio={profile.user.bio} />
        <InterestCard
          title="Favorites"
          details={"Holidays, Shoes, Food"}
          slides={favorites}
        />
        <InterestCard
          title="Food"
          details={"Pizza, Pasta, Königsberger Klopse"}
          slides={food}
        />
        <InterestCard
          title="Your Network"
          details={
            "Here you can find more details about people related to you."
          }
          slides={["http://socialengineindia.com/images/home/expert1.png"]}
        />
      </div>
      <div className="dashboardFooter">
        <Button onClick={nextPage}>Why?</Button>
      </div>
    </div>
  );
}

export default Dashboard;
