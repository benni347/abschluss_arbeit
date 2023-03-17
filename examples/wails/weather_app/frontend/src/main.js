import "./style.css";
import "./app.css";

import { Temprature } from "../wailsjs/go/main/App";

document.querySelector("#app").innerHTML = `
      <div class="result" onload="retriveTemp()" id="result">Current Temprature in Winterthur is: </div>
      <div class="input-box" id="input">
        <button class="btn" onclick="retriveTemp()">Refresh Temprature</button>
      </div>
    </div>
`;

const resultElement = document.getElementById("result");

window.retriveTemp = function () {
  try {
    Temprature().then((result) => {
      resultElement.innerText = `${result}°C`;
    })
      .catch((err) => {
        console.error(err);
      });
  } catch (err) {
    consolel.error(err);
  }
};

retriveTemp();
