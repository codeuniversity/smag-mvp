import React from "react";
import PropTypes from "prop-types";
import { Zoom } from "react-slideshow-image";

const properties = {
  duration: 5000,
  transitionDuration: 300,
  indicators: false,
  scale: 1.4,
  arrows: false
};

const Slideshow = props => {
  return (
    <div className="slideshow">
      <Zoom className="slides" {...properties}>
        {props.slides.map((imageUrl, index) => (
          <img
            key={index}
            src={imageUrl}
            alt="slide"
            className="slideshow-image"
          />
        ))}
      </Zoom>
    </div>
  );
};

Slideshow.propTypes = {
  slides: PropTypes.array
};

export default Slideshow;
