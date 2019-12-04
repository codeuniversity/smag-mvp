var typeface;

var fontSize = 18;
var streams = [];
var fadeInterval = 1;

function preload() {
  typeface = loadFont("Barlow-SemiBold.otf");
}

function setup() {
  createCanvas(windowWidth, windowHeight);
  textSize(fontSize);

  var xstream = 0;
  for (var i = 0; i <= width / fontSize; i++) {
    var stream = new Stream();
    stream.generateLetters(xstream, random(-windowHeight, windowHeight));
    streams.push(stream);
    xstream += fontSize;
  }
}

function draw() {
  background(0, 200);
  noStroke();

  streams.forEach(function(stream) {
    stream.render();
  });
}

function Letter(x, y, speed, first, opacity) {
  this.x = x;
  this.y = y;
  this.value;
  this.speed = speed;
  this.switch = round(random(10, 80));
  this.first = first;
  this.opacity = opacity;

  this.RandomLetter = function() {
    if (frameCount % this.switch === 0) {
      this.value = String.fromCharCode(0x00 + round(random(48, 90)));
    }
  };

  this.rain = function() {
    this.y = this.y >= height ? 0 : (this.y += this.speed);
  };
}

function Stream() {
  this.letters = [];
  this.totalLetters = round(random(15, 30));
  this.speed = random(2, 5);

  this.generateLetters = function(x, y) {
    var opacity = 255;
    var first = round(random(0, 4)) === 1;
    for (var i = 0; i <= this.totalLetters; i++) {
      letter = new Letter(x, y, this.speed, first, opacity);
      letter.RandomLetter();
      this.letters.push(letter);
      opacity -= 255 / this.totalLetters / fadeInterval;
      y -= fontSize;
      first = false;
    }
  };

  this.render = function() {
    this.letters.forEach(function(letter) {
      if (letter.first) {
        fill(150, 220, 255, letter.opacity);
      } else {
        fill(42, 159, 216, letter.opacity);
      }
      text(letter.value, letter.x, letter.y);
      letter.rain();
      letter.RandomLetter();
    });
  };
}
