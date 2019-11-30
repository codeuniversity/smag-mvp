import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import { Route, Link, BrowserRouter as Router, Switch } from "react-router-dom";
import App from "./App";
import { BrowserRouter } from "react-router-dom";
import Result from "./components/Result";
import Start from "./components/Start";
import Dashboard from "./pages/Dashboard";
import Notfound from "./notfound";

const root = document.getElementById("root");

ReactDOM.render(
  <BrowserRouter>
    <div>
      <Switch>
        <Route exact path="/" component={App} />
        <Route exact path="/result" component={Result} />
        <Route exact path="/dashboard" component={Dashboard} />
        <Route exact path="/start" component={Start} />
        <Route component={Notfound} />
      </Switch>
    </div>
  </BrowserRouter>,
  root
);
