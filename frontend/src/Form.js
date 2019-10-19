import React, { Component, useState } from "react";
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
          placeholder="Full name or username"
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

// class Bla extends Component{
//     state = {
//         value: ""
//     }

//     render(){
//         const {value} = this.state
//         const {onSubmit} = this.props

//         return (
//             <div className="searchForm">
//                 <form>
//                     <input type="text" id="filter" placeholder="Full name or username" value={value} onChange={this.handleSubmit}  />
//                     <input type="button" value="Submit" onClick={props.onSubmit} />
//                 </form>
//             </div>
//         )

//     }

//     handleSubmit = (e)=>{
//         e.preventDefault()
//         this.setState({value: e.target.value})
//     }
// }
