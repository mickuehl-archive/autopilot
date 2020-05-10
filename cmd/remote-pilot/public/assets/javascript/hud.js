
const RADIUS = 3;

var trackerCanvas = document.getElementById('tracker-canvas');
var trackerCtx = trackerCanvas.getContext('2d');

// HUD elments
var display_xrel = document.getElementById('display-xrel');
var display_yrel = document.getElementById('display-yrel');
var display_status = document.getElementById('display-status');

// state
var hud = {
    // local data from the browser
    width: trackerCanvas.width,
    height: trackerCanvas.height,
    width12: trackerCanvas.width / 2,
    height12: trackerCanvas.height / 2,
    x: trackerCanvas.width / 2,
    y: trackerCanvas.height / 2,
    // data from the vehicle's OBU
    steering: 0,
    throttle: 0,
    heading: 0,
    mode: "STOPPED"
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

function displayHUD(h) {
    display_status.textContent = h.mode;
    display_xrel.textContent = h.steering.toPrecision(4);
    display_yrel.textContent = h.throttle.toPrecision(4);
}

function updateHud(data) {
    // unpack values from the update event
    hud.throttle = data.th
    hud.steering = data.st
    hud.heading = data.head
    hud.mode = data.mode
    // redraw the HUD
    displayHUD(hud)
}

function resetCanvas() {
    trackerCtx.clearRect(0, 0, trackerCanvas.width, trackerCanvas.height);

    hud.x = hud.width12;
    hud.y = hud.height12;
    state.throttle = 0;
    state.steering = 0;

    drawTracker(hud.x, hud.y);
    displayHUD(hud);
}