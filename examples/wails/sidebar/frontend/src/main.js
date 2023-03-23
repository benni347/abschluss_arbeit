import "./style.css";
import "./app.css";

import { IconGenerator } from "../wailsjs/go/main/App";

document.getElementById("icongetter").innerHTML = `
  <input type="number" name="amount" value="" id="amount">
  <buton id="submit" class="btn" onclick="icon()">Get</button>
`;

window.icon = function () {
  const amount = document.getElementById("amount");
  const amountValue = amount.value;
  console.log(amountValue);
  if (amountValue === 0 || amountValue === "0" || amountValue === "") return;
  try {
    IconGenerator(amountValue).then((result) => {
      console.log(result);
    })
      .catch((err) => {
        console.error(err);
      });
  } catch (err) {
    console.error(err);
  }
};
