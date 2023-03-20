import "./style.css";
import "./app.css";

import { Greet } from "../wailsjs/go/main/App";

document.querySelector("#app").innerHTML = `
      <div class="hello"><h1>Hello User</h1></div>
      <div class="input-box" id="input">
        <input class="input" id="name" type="text" autocomplete="off" />
        <button id="submit" class="btn" onclick="greet()">Greet</button>
      </div>
    </div>
`;

const nameElement = document.getElementById("name");
nameElement.focus();

// Setup the greet function
window.greet = function () {
  // Get name
  const name = nameElement.value;
  
  // Check if the input is empty
  if (name === "") return;

  // Call App.Greet(name)
  try {
    Greet(name)
      .then(() => {
        // Update result with data back from App.Greet()
        nameElement.value = "";
        nameElement.focus();
      })
      .catch((err) => {
        console.error(err);
      });
  } catch (err) {
    console.error(err);
  }
};
