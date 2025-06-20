const protocol = window.location.protocol === "https:" ? "wss" : "ws";
const host = window.location.host;
let ws;
let pingIntervalId;
let pongTimeoutId;


function sendCurrentPath() {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      action: "setLocation",
      path: window.location.pathname,
    }));
  }
}

function connect() {
  ws = new WebSocket(`${protocol}://${host}/ws?path=${encodeURIComponent(window.location.pathname)}`);

  ws.onopen = () => {
     sendCurrentPath();
    window.addEventListener("popstate", sendCurrentPath);
    console.log("WebSocket connected");
  };

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      console.log("WebSocket message received:", data);

      if (data.type === "pong") {
        console.log("Received app-level pong (if you keep app pings)");
      }

      // your existing update logic
      if (data.action === "update" && data.id && data.html) {
        const target = document.getElementById(data.id);
        if (target) {
          target.innerHTML = data.html;
        } else {
          console.warn(`Element with id "${data.id}" not found.`);
        }
      }
    } catch (e) {
      console.error("Failed to handle WebSocket message:", e);
    }
  };

  ws.onclose = () => {
    console.log("WebSocket connection closed");
    clearInterval(pingIntervalId);
    clearTimeout(pongTimeoutId);
    setTimeout(connect, 5_000);
  };

  ws.onerror = (err) => {
    console.error("WebSocket error:", err);
    ws.close();
  };

  router.afterEach((to) => {
    sendCurrentPath(to.fullPath);
  });

}

connect();

