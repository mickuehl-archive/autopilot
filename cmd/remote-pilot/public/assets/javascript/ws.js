
var WSREADY = false;

var ws = new WebSocket("ws://localhost:3000/hud");

ws.onopen = function () {
    WSREADY = true;
    console.log("socket connection is ready")
};

ws.onmessage = function (evt) {
    var received_msg = evt.data;
    console.log("->" + received_msg)
};

ws.onclose = function () {
    console.log("socket connection closed")
    WSREADY = false;
};

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
