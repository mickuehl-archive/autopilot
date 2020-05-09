
const BORDER = 6;

var animation;

// state
var state = {
    mode: "STOPPED",
    throttle: 0,
    steering: 0
}

function updatePosition(e) {
    hud.x += e.movementX;
    hud.y += e.movementY;
    if (hud.x > hud.width - (BORDER / 2)) {
        hud.x = hud.width - (BORDER / 2);
    }
    if (hud.y > hud.height - (BORDER / 2)) {
        hud.y = hud.height - (BORDER / 2);
    }
    if (hud.x < BORDER) {
        hud.x = BORDER;
    }
    if (hud.y < BORDER) {
        hud.y = BORDER;
    }

    state.steering = (hud.x - hud.width12) / hud.width12;
    state.throttle = -1 * ((hud.y - hud.height12) / hud.height12);

    sendState(state);
    displayHUD(state);

    if (!animation) {
        animation = requestAnimationFrame(function () {
            animation = null;
            drawTracker(hud.x, hud.y);
        });
    }
}

function lockChangeAlert() {
    if (document.pointerLockElement === trackerCanvas || document.mozPointerLockElement === trackerCanvas) {
        state.mode = "DRIVING";
        document.addEventListener("mousemove", updatePosition, false);
    } else {
        state.mode = "STOPPED";
        document.removeEventListener("mousemove", updatePosition, false);
        resetCanvas();
        sendState(state);
    }
}

// draw the canvas for the first time
resetCanvas();

// pointer lock object forking for cross browser

trackerCanvas.requestPointerLock = trackerCanvas.requestPointerLock ||
    trackerCanvas.mozRequestPointerLock;

document.exitPointerLock = document.exitPointerLock ||
    document.mozExitPointerLock;

trackerCanvas.onclick = function () {
    trackerCanvas.requestPointerLock();
};

// pointer lock event listeners

// Hook pointer lock state change events for different browsers
document.addEventListener('pointerlockchange', lockChangeAlert, false);
document.addEventListener('mozpointerlockchange', lockChangeAlert, false);
