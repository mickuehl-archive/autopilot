// initialize stuff

const RADIUS = 3;
const BORDER = 5;

var animation;

var canvas = document.getElementById('tracker-canvas');
var ctx = canvas.getContext('2d');

// HUD elments
var display_xrel = document.getElementById('display-xrel');
var display_yrel = document.getElementById('display-yrel');
var display_status = document.getElementById('display-status');

var x = canvas.width / 2;
var y = canvas.height / 2;

var rel_width = canvas.width / 2;
var rel_height = canvas.height / 2;
var xrel = 0;
var yrel = 0;

var status = "STOPPED"

function degToRad(degrees) {
    var result = Math.PI / 180 * degrees;
    return result;
}

function drawCanvas() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.fillStyle = "#f00";
    ctx.beginPath();
    ctx.arc(x, y, RADIUS, 0, degToRad(360), true);
    ctx.fill();
}

function resetCanvas() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    x = canvas.width / 2;
    y = canvas.height / 2;
    xrel = 0;
    yrel = 0;
    drawCanvas();
    displayHUD();
}

function displayHUD() {
    display_status.textContent = status
    display_xrel.textContent = xrel.toPrecision(4);
    display_yrel.textContent = yrel.toPrecision(4);
}

function updatePosition(e) {
    x += e.movementX;
    y += e.movementY;
    if (x > canvas.width - BORDER) {
        x = canvas.width - BORDER;
    }
    if (y > canvas.height - BORDER) {
        y = canvas.height - BORDER;
    }
    if (x < BORDER) {
        x = BORDER;
    }
    if (y < BORDER) {
        y = BORDER;
    }

    xrel = (x - rel_width) / rel_width;
    yrel = -1 * ((y - rel_height) / rel_height);
    displayHUD();

    if (!animation) {
        animation = requestAnimationFrame(function () {
            animation = null;
            drawCanvas();
        });
    }
}

function lockChangeAlert() {
    if (document.pointerLockElement === canvas || document.mozPointerLockElement === canvas) {
        status = "DRIVING";
        document.addEventListener("mousemove", updatePosition, false);
    } else {
        status = "STOPPED";
        document.removeEventListener("mousemove", updatePosition, false);
        resetCanvas();
    }
}

// draw the canvas for the first time
resetCanvas();

// pointer lock object forking for cross browser

canvas.requestPointerLock = canvas.requestPointerLock ||
    canvas.mozRequestPointerLock;

document.exitPointerLock = document.exitPointerLock ||
    document.mozExitPointerLock;

canvas.onclick = function () {
    canvas.requestPointerLock();
};

// pointer lock event listeners

// Hook pointer lock state change events for different browsers
document.addEventListener('pointerlockchange', lockChangeAlert, false);
document.addEventListener('mozpointerlockchange', lockChangeAlert, false);
