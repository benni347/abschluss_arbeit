import { Publisher, Dial } from "../wailsjs/go/main/App";

(() => {
  // expectingMessage is set to true
  // if the user has just submitted a message
  // and so we should scroll the next message into view when received.
  let expectingMessage = false;
  function dial() {
    let msg = "";
    let error = "";
    console.info(`location.host Type: ${typeof location.host}`);
    msg, (error = Dial(location.host));
    console.log(msg);
    if (msg === "") {
      appendLog(`WebSocket Disconnected code: ${error}`, true);
    } else {
      const p = appendLog(msg);
      if (expectingMessage) {
        p.scrollIntoView();
        expectingMessage = false;
      }
    }
    /* 
    console.log("dialing")
    const conn = new WebSocket(`ws://${location.host}/ws`)
    conn.addEventListener("close", ev => {
      appendLog(`WebSocket Disconnected code: ${ev.code}, reason: ${ev.reason}`, true)
      if (ev.code !== 1001) {
        appendLog("Reconnecting in 1s", true)
        setTimeout(dial, 1000)
      }
    })
    conn.addEventListener("open", ev => {
      console.info("websocket connected")
    })

    // This is where we handle messages received.
    conn.addEventListener("message", ev => {
      if (typeof ev.data !== "string") {
        console.error("unexpected message type", typeof ev.data)
        return
      }
      const p = appendLog(ev.data)
      if (expectingMessage) {
        p.scrollIntoView()
        expectingMessage = false
      }
    })
    */
  }
  const messageLog = document.getElementById("message-log");
  const publishForm = document.getElementById("publish-form");
  const messageInput = document.getElementById("message-input");

  dial();

  // appendLog appends the passed text to messageLog.
  function appendLog(text, error) {
    const p = document.createElement("p");
    // Adding a timestamp to each message makes the log easier to read.
    p.innerText = `${new Date().toLocaleTimeString()}: ${text}`;
    if (error) {
      p.style.color = "red";
      p.style.fontStyle = "bold";
    }
    messageLog.append(p);
    return p;
  }

  // onsubmit publishes the message from the user when the form is submitted.
  publishForm.onsubmit = async (ev) => {
    ev.preventDefault();

    const msg = messageInput.value;
    if (msg === "") {
      return;
    }
    messageInput.value = "";

    expectingMessage = true;
    let error = Publisher(msg, location.host);
    appendLog(`Publish failed: ${msg}`, error);
  };
})();

console.log(Publisher("hello", location.host));
