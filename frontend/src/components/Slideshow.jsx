import React, { Component } from "react";
import { Zoom } from "react-slideshow-image";

const properties = {
  duration: 5000,
  transitionDuration: 300,
  indicators: false,
  scale: 1.4,
  arrows: false
};

class Slideshow extends Component {
  render() {
    return (
      <div className="slideshow">
        <Zoom className="slides" {...properties}>
          {this.props.slides.map((imageUrl, index) => (
            <img key={index} src={imageUrl} alt="slide" />
          ))}
        </Zoom>
      </div>
    );
  }
}

export default Slideshow;
