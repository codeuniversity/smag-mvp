import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import { Route, Link, BrowserRouter as Router, Switch } from "react-router-dom";
import App from "./App";
import { BrowserRouter } from "react-router-dom";
import Result from "./components/Result";
import Dashboard from "./pages/Dashboard";
import Notfound from "./notfound";
import FlowWrapper from "./components/FlowWrapper";
import Popup from "./components/Popup";
import Greeting from "./pages/Greeting.jsx";
import Network from "./pages/Network.jsx";
import Endscreen from "./pages/endscreen";
import SearchProfile from "./pages/SearchProfile";
import ExampleProfileSelection from "./pages/ExampleProfileSelection";

const root = document.getElementById("root");

ReactDOM.render(
  <BrowserRouter>
    <div>
      <Switch>
        <Route exact path="/" component={FlowWrapper} />
        <Route exact path="/result" component={Result} />
        <Route exact path="/greeting" component={Greeting} />
        <Route exact path="/popup" component={Popup} />
        <Route exact path="/endscreen" component={Endscreen} />
        <Route exact path="/dashboard" component={Dashboard} />
        <Route exact path="/app" component={App} />
        <Route exact path="/network" component={Network} />
        <Route exact path="/search-profile" component={SearchProfile} />
        <Route
          exact
          path="/example-profile-selection"
          component={ExampleProfileSelection}
        />
        <Route component={Notfound} />
      </Switch>
    </div>
  </BrowserRouter>,
  root
);
