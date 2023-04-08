import "./style.css";
import "./app.css";

import { DeleteFile, IconGenerator } from "../wailsjs/go/main/App";

document.getElementById("icongetter").innerHTML = `
  <input type="number" name="amount" value="" id="amount">
  <buton id="submit" class="btn" onclick="icon()">Get</button>
`;

const iconPreviewDiv = document.getElementById("icon_preview");

window.icon = function () {
  const amount = document.getElementById("amount");
  const amountValue = amount.value;
  console.log(amountValue);
  if (amountValue === 0 || amountValue === "0" || amountValue === "") return;
  try {
    IconGenerator(amountValue).then((result) => {
      for (let i = 0; i < result.length; i++) {
        const fileName = result[i];
        const imgElement = document.createElement("img");
        const filePath = "./src/assets/images/" + fileName;
        imgElement.src = filePath;

        // Add the image element to the icon preview div
        iconPreviewDiv.appendChild(imgElement);
      }
    })
      .catch((err) => {
        console.error(err);
      });
  } catch (err) {
    console.error(err);
  }
  try {
    DeleteFile(fileName);
  } catch (err) {
    console.error(err);
  }
};
