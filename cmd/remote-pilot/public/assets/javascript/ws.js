
// flag to controll sending of data
var WSREADY = false;

// open a connection to the vehicle
var ws = new WebSocket("ws://localhost:3000/hud");

// Send the local state of the pilot to the vehicle. 
// The vehicle SHOULD reply with an update message with actual values to be 
// display in the HUD / overlay
function sendState(s) {
    if (WSREADY == true) {
        var data = {
            th: s.throttle,
            st: s.steering,
            mode: s.mode,
            ts: Date.now()
        }
        ws.send(JSON.stringify(data))
    }
}

// callback to update the hud state based on the data sent by the vehicle
ws.onmessage = function (event) {
    // unmarshal the update event and update the HUD
    updateHud(JSON.parse(event.data).Data);
};

// websocket callbacks used to maintain the connection state
ws.onopen = function () {
    WSREADY = true;
    console.log("socket connection is ready")
};

ws.onclose = function () {
    console.log("socket connection closed")
    WSREADY = false;
};


