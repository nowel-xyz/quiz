const ws = new WebSocket("ws://localhost:3000/ws");

ws.onopen = () => {
  console.log("WebSocket connected");
};

ws.onmessage = (event) => {
  try {
    const data = JSON.parse(event.data);

    if (data.type === "update" && data.id && data.html) {
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
};
