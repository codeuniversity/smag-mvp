import React, { Component, useState } from "react";
import { withRouter, history } from "react-router";

function IGPost({ post }) {
  return (
    <div className="ig-post">
      <div className="ig-post-padding">
        <a target="_blank" href={"https://instagram.com/p/" + post.shortcode}>
          <img className="post-image" src={post.img} />
        </a>
      </div>
    </div>
  );
}

export default IGPost;
