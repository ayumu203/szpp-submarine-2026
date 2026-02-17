const BOARD_SIZE = 5;
const fieldTable = document.querySelector("table#field");

/**
 * @description jsonファイルからモックデータを取得する
 * @returns {object} モックデータ
 */
async function getMock() {

      const response = await fetch("mock/mock.json");

      if (!response.ok) {
            throw new Error("JSONデータの読み込みに失敗しました");
      }
      const data = await response.json();
      return data;
}


/**
 * @description 潜水艦の座標一覧を取得する
 * @returns {{x: number, y: number}[]} 各潜水艦の座標
 */
async function getSubmarinePositions() {
      const data = await getMock();
      const submarines = Object.values(data.State.GetGameStateResponse.allyBoard.submarines);

      return submarines.map((submarine) => ({
            x: submarine.x,
            y: submarine.y,
      }));
}

/** @description 表の内容をすべて空にする */
function clearSubmarines() {
      for (let row = 1; row <= BOARD_SIZE; row++) {
            for (let col = 1; col <= BOARD_SIZE; col++) {
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
      if (!fieldTable) {
            throw new Error("潜水艦を表示するテーブルが存在しません");
      }

      const positions = await getSubmarinePositions();
      clearSubmarines();

      positions.forEach(({ x, y }) => {
            if (x < 1 || x > BOARD_SIZE || y < 1 || y > BOARD_SIZE) return;
            fieldTable.rows[y].cells[x].textContent = "●";
      });
}

document.addEventListener("DOMContentLoaded", renderSubmarines);
