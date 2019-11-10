import React, { Component } from "react";
import App, { users, userdata, user } from "../App";
import { withRouter } from "react-router";
import IGlogo from "./img/instagram.png";
import FBlogo from "./img/fb.png";
import Twitterlogo from "./img/twitter.png";
import LIlogo from "./img/linkedin.png";
import IGPost from "./IGPost";

import {
  UserSearchServiceClient,
  UserIdRequest
} from "../protofiles/usersearch_grpc_web_pb";

export class Result extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      user: this.props.location.state && this.props.location.state.user,
      posts: []
    };
  }

  componentDidMount() {
    const user = this.state.user;
    if (!user) {
      console.log(":(");
      return;
    }
    const client = new UserSearchServiceClient("http://localhost:4000");
    const userIdRequest = new UserIdRequest();
    userIdRequest.setUserId(user.id);
    client.getInstaPostsWithUserId(userIdRequest, {}, (err, response) => {
      if (err) {
        console.log(err);
        return;
      }
      console.log(response);
      const posts = response.getInstaPostsList();
      const postdata = posts
        .map(posts => ({
          id: posts.getPostId(),
          img: posts.getImgUrl(),
          caption: posts.getCaption(),
          shortcode: posts.getShortCode()
        }))
        .filter(post => !!post.img);
      // this.props.history.push({
      //   pathname: "/results",
      //   state: { results: postdata }
      // });
      console.log(postdata);
      this.setState({
        posts: postdata
      });
    });
  }

  render() {
    const { user, posts } = this.state;
    if (!user) {
      return "no user found";
    }

    return (
      <div className="body">
        <div className="container-card">
          <div className="column-one-third"></div>
          <div className="column-one-third">
            <div className="box">
              <div>
                <img className="avatar-image" src={user.avatarurl} />
              </div>
              <div className="headline">Hello {user.realname}</div>
              <div className="sub-headline" key={user.username}>
                <a
                  target="_blank"
                  href={"https://instagram.com/" + user.username}
                >
                  @{user.username}
                </a>
              </div>
              <div className="body-text">{user.bio}</div>
            </div>
          </div>
          <div className="column-one-third"></div>
        </div>

        <div className="container-plattforms">
          <div className="column-one-sixed"></div>
          <div className="column-four-sixed">
            <div className="social-icon-box">
              <div className="sub-headline">
                Your public Social Media Accounts:
              </div>
              <div>
                <a
                  target="_blank"
                  href={"https://instagram.com/" + user.username}
                >
                  <img className="social-icon" src={IGlogo} />
                </a>
              </div>
              <div>
                <img className="social-icon" src={FBlogo} />
              </div>
              <div>
                <img className="social-icon" src={Twitterlogo} />
              </div>
              <div>
                <img className="social-icon" src={LIlogo} />
              </div>
            </div>
          </div>
          <div className="column-one-sixed"></div>
        </div>

        <div className="container-plattforms">
          <div classNme="column-one-sixed"></div>
          <div className="column-four-sixed-posts">
            {posts.map(post => (
              <IGPost key={post.shortcode} post={post} />
            ))}
          </div>
          <div classNme="column-one-sixed"></div>
        </div>
      </div>
    );
  }
}
export default withRouter(Result);