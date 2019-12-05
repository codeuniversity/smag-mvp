import React, { Component } from "react";
import Button from "../components/Button";
import "./../index.css";
import H1 from "../components/H1";
import H2 from "../components/H2";

function GroupIntent({ nextPage }) {
  return (
    <div className="container-center">
      <div className="column-center">
        <div className="greeting">
          <H1>Anyone can do this</H1>
          <H2>We had limited time and money.</H2>
          <H2>Imagine what others could do with this power.</H2>
          <div style={{ marginTop: 50 }}>
            <Button onClick={nextPage}>Noted.</Button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default GroupIntent;
