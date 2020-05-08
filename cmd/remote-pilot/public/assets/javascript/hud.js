
const RADIUS = 3;

var trackerCanvas = document.getElementById('tracker-canvas');
var trackerCtx = trackerCanvas.getContext('2d');

// HUD elments
var display_xrel = document.getElementById('display-xrel');
var display_yrel = document.getElementById('display-yrel');
var display_status = document.getElementById('display-status');

// state
var hud = {
    width: trackerCanvas.width,
    height: trackerCanvas.height,
    width12: trackerCanvas.width / 2,
    height12: trackerCanvas.height / 2,
    x: trackerCanvas.width / 2,
    y: trackerCanvas.height / 2
}

function degToRad(degrees) {
    var result = Math.PI / 180 * degrees;
    return result;
}

function drawTracker(x, y) {
    trackerCtx.clearRect(0, RADIUS, trackerCanvas.width, trackerCanvas.height);
    trackerCtx.fillStyle = "#f00";
    trackerCtx.beginPath();
    trackerCtx.arc(x, y, RADIUS, 0, degToRad(360), true);
    trackerCtx.fill();
}

function displayHUD(s) {
    display_status.textContent = s.mode;
    display_xrel.textContent = s.steering.toPrecision(4);
    display_yrel.textContent = s.throttle.toPrecision(4);
}

function resetCanvas() {
    trackerCtx.clearRect(0, 0, trackerCanvas.width, trackerCanvas.height);

    hud.x = hud.width12;
    hud.y = hud.height12;
    state.throttle = 0;
    state.steering = 0;

    drawTracker(hud.x, hud.y);
    displayHUD(state);
}