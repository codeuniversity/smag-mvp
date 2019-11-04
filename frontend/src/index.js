import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import { Route, Link, BrowserRouter as Router, Switch } from "react-router-dom";
import App from "./App";
import { BrowserRouter } from "react-router-dom";
import Results from "./components/Results";
import Notfound from "./notfound";

const root = document.getElementById("root");

ReactDOM.render(
  <BrowserRouter>
    <div>
      <Switch>
        <Route exact path="/" component={App} />
        <Route exact path="/results" component={Results} />
        <Route component={Notfound} />
      </Switch>
    </div>
  </BrowserRouter>,
  root
);
