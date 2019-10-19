import React from 'react'
import ReactDOM from 'react-dom'
import './index.css'
import { Route, Link, BrowserRouter as Router } from 'react-router-dom'
import App from './App'
import Results from './Results'
const routing = (
  <Router>
    <div>
      <Route exact path="/" component={App} />
      <Route path="/results" component={Results} />
    </div>
  </Router>
)
ReactDOM.render(routing, document.getElementById('root'))
