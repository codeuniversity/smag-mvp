import React, { useState, useEffect } from "react";
import "./FaceHitAnimation.css";

function FaceHitAnimation(props) {
  // hack to freeze faceHits on first mount
  const [faceHits, setFaceHits] = useState({});
  useEffect(() => {
    setFaceHits(props.faceHits);
  }, []);

  const images = Object.entries(faceHits)
    .map(([postId, faces]) => ({
      postId,
      weight: faces.length,
      imageSrc: faces[0].fullImageSrc
    }))
    .sort((a, b) => b.weigth - a.weight)
    .reverse();

  const [imageCoordinates, setImageCoordinates] = useState([]);

  useEffect(() => {
    if (imageCoordinates.length >= images.length) {
      const timeoutId = setTimeout(props.onAnimationFinished, 2000);
      return () => clearTimeout(timeoutId);
    }

    const intervalId = setInterval(() => {
      setImageCoordinates(prevCoordinates => [
        ...prevCoordinates,
        { top: randomCoordinate(), left: randomCoordinate() }
      ]);
    }, 250);
    return () => clearInterval(intervalId);
  }, [imageCoordinates.length, images.length]);

  return (
    <>
      {imageCoordinates.map((coord, index) => (
        <BackgroundImage
          key={images[index].postId}
          imageSrc={images[index].imageSrc}
          top={coord.top}
          left={coord.left}
          last={index === images.length - 1}
        />
      ))}
    </>
  );
}

export default FaceHitAnimation;

function BackgroundImage({ imageSrc, top, left, last }) {
  if (last) {
    return <img className="BackgroundImage-center" src={imageSrc} />;
  }

  return (
    <img
      className="BackgroundImage"
      src={imageSrc}
      style={{ top: `calc(${top}% - 300px)`, left: `calc(${left}% - 300px)` }}
    />
  );
}

function randomCoordinate() {
  return Math.round(Math.random() * 100);
}
