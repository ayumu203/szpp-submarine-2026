const table = document.querySelector("table#field");

async function getMock() {

      const response = await fetch("mock/mock.json");
      const data = await response.json();

      return data;
}

async function getPos() {
      const data = await getMock();
      const submarines = Object.values(data.State.GetGameStateResponse.allyBoard.submarines);
      const xPos = [];
      const yPos = [];

      submarines.forEach(submarine => {
            xPos.push(submarine.x);
            yPos.push(submarine.y);
      });

      clearSubmarines();
      renderSubmarines(xPos, yPos);
}

function clearSubmarines() {
      for (let row = 1; row <= 5; row++) {
            for (let col = 1; col <= 5; col++) {
                  table.rows[row].cells[col].textContent = "";
            }
      }
}

function renderSubmarines(xPos, yPos) {
      for (let i = 0; i < 4; i++) {
            x = xPos[i];
            y = yPos[i];
            if (x < 1 || x > 5 || y < 1 || y > 5) return;
            table.rows[y].cells[x].textContent = "●";
      }
}

document.addEventListener("DOMContentLoaded", getPos);
