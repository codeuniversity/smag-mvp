import React, { Component } from "react";

export class CameraFeed extends Component {
  /**
   * Processes available devices and identifies one by the label
   * @memberof CameraFeed
   * @instance
   */
  processDevices(devices) {
    devices.forEach(device => {
      this.setDevice(device);
    });
  }

  /**
   * Sets the active device and starts playing the feed
   * @memberof CameraFeed
   * @instance
   */
  async setDevice(device) {
    const { deviceId } = device;
    const stream = await navigator.mediaDevices.getUserMedia({
      audio: false,
      video: { deviceId }
    });
    this.videoPlayer.srcObject = stream;
    this.videoPlayer.play();
  }

  /**
   * On mount, grab the users connected devices and process them
   * @memberof CameraFeed
   * @instance
   * @override
   */
  async componentDidMount() {
    const cameras = await navigator.mediaDevices.enumerateDevices();
    this.processDevices(cameras);

    this.intervalId = setInterval(() => {
      this.takePhoto();
    }, 2000);
  }

  componentWillUnmount() {
    clearInterval(this.intervalId);
  }

  /**
   * Handles taking a still image from the video feed on the camera
   * @memberof CameraFeed
   * @instance
   */
  takePhoto = () => {
    const { onFileSubmit } = this.props;
    const context = this.canvas.getContext("2d");
    context.drawImage(this.videoPlayer, 0, 0, 800, 600);
    this.canvas.toBlob(onFileSubmit);
  };

  stopTakingPictures = () => {
    clearInterval(this.intervalId);
  };

  render() {
    return (
      <div className="Button-center">
        <video ref={ref => (this.videoPlayer = ref)} width="800" heigh="600" />
        <div className="c-camera-feed__stage">
          <canvas width="800" height="600" ref={ref => (this.canvas = ref)} />
        </div>
      </div>
    );
  }
}
