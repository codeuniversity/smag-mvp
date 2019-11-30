import React, { Component } from "react";
import { Slide } from "react-slideshow-image";

const properties = {
  duration: 3000,
  transitionDuration: 400,
  infinite: true,
  indicators: false,
  arrows: false
};

class Slideshow extends Component {
  render() {
    return (
      <div className="slideshow">
        <Slide className="slides" {...properties}>
          {this.props.slides.map(imageUrl => (
            <div className="each-slide">
              <img src={imageUrl} alt="slide" />
            </div>
          ))}
        </Slide>
      </div>
    );
  }
}

export default Slideshow;
