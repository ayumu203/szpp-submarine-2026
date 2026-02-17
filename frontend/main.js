const table = document.querySelector("table#field");

/**
 * @description jsonファイルからモックデータを取得する
 * @returns {object} モックデータ
 */
async function getMock() {

      const response = await fetch("mock/mock.json");
      const data = await response.json();

      return data;
}

/** @description 表の内容をすべて空にする */
function clearSubmarines() {
      for (let row = 1; row <= 5; row++) {
            for (let col = 1; col <= 5; col++) {
                  if (table) {
                        table.rows[row].cells[col].textContent = "";
                  }
            }
      }
}

/**
 * @description 潜水艦の位置を●で表示する
 * @param {number[]} xPos x座標のリスト
 * @param {number[]} yPos y座標のリスト
*/
function renderSubmarines(xPos, yPos) {
      for (let i = 0; i < xPos.length; i++) {
            const x = xPos[i];
            const y = yPos[i];
            if (x < 1 || x > 5 || y < 1 || y > 5) continue;
            table.rows[y].cells[x].textContent = "●";
      }
}

/**
 * @description 潜水艦の位置を取得、表示する
*/
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


document.addEventListener("DOMContentLoaded", getPos);
