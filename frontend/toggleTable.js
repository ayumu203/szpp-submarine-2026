// 「予測表示切り替え」ボタンを押したときに、「味方潜水艦」「予測度」の表を切り替える処理
const DISPLAY_MODES = ["ally", "prediction"]

/**
 * @description 表示するテーブルを切り替える
 * @param {string} currentMode 今表示しているテーブルの状態
 * @returns {string} 次の状態
 */
function toggleDisplayMode(currentMode) {
  return currentMode === "ally" ? "prediction" : "ally";
}

function clearBoardView() {
  const $rows = $("#field tr");
  for (let y = 1; y <= BOARD_SIZE; y++) {
    for (let x = 1; x <= BOARD_SIZE; x++) {
      $rows.eq(y).children("td").eq(x).text("").removeClass("prediction-cell");
    }
  }
}

/**
 * @description 潜水艦がいるセルをmarkerで表示する
 * @param {object} submarines すべての潜水艦
 * @param {string} marker どうやってUIに表示するか(味方の場合は"●")
 */
function renderBoardBySubmarines(submarines, marker) {
  const submarines_list = Object.values(submarines);

  submarines_list.forEach((submarine) => {
    if (submarine.x < 1 || submarine.x > BOARD_SIZE || submarine.y < 1 || submarine.y > BOARD_SIZE) return;

    if (submarine.sunk) {
      fieldTable.rows[submarine.y].cells[submarine.x].textContent = "S";
    } else {
      fieldTable.rows[submarine.y].cells[submarine.x].textContent = marker;
    }
  });
}

/**
 * @description 敵の位置の予測表示を描写する
 * @param {object} predictionBoard jsonから取得した、敵の存在確率(?)
 */
function renderPredictionBoard(predictionBoard) {
  
}

/**
 * @description uiState.displayModeに応じてテーブルを表示する
 * @returns 謎
 */
async function renderDisplayMode() {
  const data = await getMock();
  const state = data.State.GetGameStateResponse;

  clearBoardView();

  if (uiState.displayMode === "ally") {
    renderBoardBySubmarines(state.allyBoard.submarines, "●");
    return;
  }

  if (uiState.displayMode === "prediction") {
    renderPredictionBoard(state.predictionBoard);
  }
}

function bindDisplayToggle() {
  $("#btn-display")
    .off("click.displayMode")
    .on("click.displayMode", async function () {
      uiState.displayMode = toggleDisplayMode(uiState.displayMode);
      console.log(uiState.displayMode);
      await renderDisplayMode();
    });
}