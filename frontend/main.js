const fieldTable = document.querySelector("table#field");

/**
 * @description jsonファイルからモックデータを取得する
 * @returns {object} モックデータ
 */
async function getMock() {

      const response = await fetch("mock/mock.json");
      const data = await response.json();

      return data;
}


/**
 * @description 潜水艦のx座標を取得する
 * @returns {number[]} 各潜水艦のx座標のリスト
*/
async function getPosX() {
      const data = await getMock();
      const submarines = Object.values(data.State.GetGameStateResponse.allyBoard.submarines);
      const xPos = [];
      
      submarines.forEach(submarine => {
            xPos.push(submarine.x);
      });
      return xPos;
}

/**
 * @description 潜水艦のy座標を取得する
 * @returns {number[]} 各潜水艦のy座標のリスト
*/
async function getPosY() {
      const data = await getMock();
      const submarines = Object.values(data.State.GetGameStateResponse.allyBoard.submarines);
      const yPos = [];
      
      submarines.forEach(submarine => {
            yPos.push(submarine.y);
      });
      return yPos;
}

/** @description 表の内容をすべて空にする */
function clearSubmarines() {
      for (let row = 1; row <= 5; row++) {
            for (let col = 1; col <= 5; col++) {
                  if (fieldTable) {
                        fieldTable.rows[row].cells[col].textContent = "";
                  }
            }
      }
}

/**
 * @description 潜水艦の位置を●で表示する
*/
async function renderSubmarines() {
      const xPos = await getPosX();
      const yPos = await getPosY();

      clearSubmarines();

      for (let i = 0; i < xPos.length; i++) {
            const x = xPos[i];
            const y = yPos[i];
            if (x < 1 || x > 5 || y < 1 || y > 5) continue;
            fieldTable.rows[y].cells[x].textContent = "●";
      }
}

document.addEventListener("DOMContentLoaded", renderSubmarines);
