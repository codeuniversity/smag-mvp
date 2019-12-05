import React, { Component, useState } from "react";
import { withRouter, history } from "react-router";

// import {onSubmit} from './App';

// eslint-disable-next-line
function Form(props) {
  const [value, setValue] = useState("");

  return (
    <div className="searchForm">
      <form
        onSubmit={e => {
          e.preventDefault();
          props.onSubmit(value);
        }}
      >
        <input
          type="text"
          id="filter"
          placeholder="Username"
          value={value}
          onChange={e => {
            e.preventDefault();
            setValue(e.target.value);
          }}
        />
        <button>Submit</button>
      </form>
    </div>
  );
}

export default Form;
