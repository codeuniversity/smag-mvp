import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import { Route, Link, BrowserRouter as Router, Switch } from "react-router-dom";
import App from "./App";
import { BrowserRouter } from "react-router-dom";
import Result from "./components/Result";
import Notfound from "./notfound";
import FlowWrapper from "./components/FlowWrapper";
import Popup from "./components/Popup";
import Greeting from "./pages/Greeting.jsx";

const root = document.getElementById("root");

ReactDOM.render(
  <BrowserRouter>
    <div>
      <Switch>
        <Route exact path="/" component={App} />
        <Route exact path="/result" component={Result} />
        <Route exact path="/greeting" component={Greeting} />
        <Route exact path="/popup" component={Popup} />
        <Route exact path="/start" component={FlowWrapper} />
        <Route component={Notfound} />
      </Switch>
    </div>
  </BrowserRouter>,
  root
);
