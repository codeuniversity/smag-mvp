import React, { Component } from "react";
import App, { users, userdata, user } from "../App";
import { withRouter } from "react-router";
import IGlogo from "./img/instagram.png";
import FBlogo from "./img/fb.png";
import Twitterlogo from "./img/twitter.png";
import LIlogo from "./img/linkedin.png";

export class Results extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      users:
        (this.props.location.state && this.props.location.state.results) || []
    };
  }

  render() {
    const results = this.state.users.map(function(result) {
      return (
        <div className="body">
          <div className="container-card">
            <div className="column-one-third"></div>
            <div className="column-one-third">
              <div className="box">
                <div>
                  <img className="avatar-image" src={result.avatarurl} />
                </div>
                <div className="headline">Hello {result.realname}</div>
                <div className="sub-headline" key={result.username}>
                  <a
                    target="_blank"
                    href={"https://instagram.com/" + result.username}
                  >
                    @{result.username}
                  </a>
                </div>
                <div className="body-text">{result.bio}</div>
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
                    href={"https://instagram.com/" + result.username}
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
        </div>
      );
    });
    return <div>{results}</div>;
  }
}
export default withRouter(Results);
/*function generateTableHead(table, data) {
      let thead = table.createTHead();
      let row = thead.insertRow();
      for (let key of data) {
        let th = document.createElement("th");
        let text = document.createTextNode(key);
        th.appendChild(text);
        row.appendChild(th);
      }
    }
    function generateTable(table, data) {
      for (let element of data) {
        let row = table.insertRow();
        for (key in element) {
          let cell = row.insertCell();
          let text = document.createTextNode(element[key]);
          cell.appendChild(text);
        }
      }
    }
    let table = document.querySelector("table");
    let data = Object.keys(userdata[0]);
    generateTableHead(table, data);
    generateTable(table, userdata);*/
