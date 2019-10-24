import React, { Component } from "react";
import App, { users, userdata, user } from "../App";

export class Results extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      users: []
    };
  }

  render() {
    console.log(this.props.results);
    //return <div>Test</div>;
    const results = this.props.results.map(function(result) {
      return (
        <div>
          {result.realname}
          {result.bio}
          <img src={result.avatarurl} />
        </div>
      );
    });
    return <div>{results}</div>;
  }
}
export default Results;
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
