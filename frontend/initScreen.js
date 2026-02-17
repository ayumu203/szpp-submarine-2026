// ゲーム開始時に、画面の初期化を行う処理を記述する
const uiState = {
    mode: "idle",  // idle | attack | moveの3種類
    displayMode: "ally", // ally | enemy | prediction
    selectedCell: null
};

/**
 * @description 潜水艦を表示するテーブルが存在するかを確認する
 * @returns {boolean} true→存在する false→存在しない
 */
function validateRequiredElements() {
    const fieldTable = document.querySelector("table#field");
    if (!fieldTable) {
        console.error("潜水艦を表示するテーブルが存在しません");
        return false;
    }
    return true;
}

/**
 * @description uiStateを初期状態に戻す
 */
function resetUiState() {
    uiState.mode = "idle";
    uiState.displayMode = "ally";
    uiState.selectedCell = null;
}

/**
 * @description 各ボタンを押したときにuiStateを更新する
 */
function changeUiStateByClick() {

    // 「攻撃」ボタンを押したら、uiStateのmodeをattackに変更する
    $("#btn-attack").off("click").on("click", () => {
        uiState.mode = "attack";
        console.log(uiState.mode);
    })

    // 「移動」ボタンを押したら、uiStateのmodeをmoveに変更する
    // 頑張って書く


    // 「予測表示切り替え」ボタンを押したら、uiStateのdisplayModeを"適切に"変更する
    // 頑張って書く


    // 「戻る」ボタンを押したら、uiStateのselectedCellをnullに、modeをidleに変更する
    // 頑張って書く


    // 「決定」ボタンを押したら、現在のuiStateを出力(console.log)する
    // 頑張って書く

}

/**
 * @description ゲーム開始時の画面を作成する
 */
async function initializeScreen() {
    // テーブルが存在しなかったらreturn

    // UIをリセットする関数を呼び出す

    // uiStateを変更する関数を呼び出す

    // main.jsの関数を使って、潜水艦を表示する

}

$(initializeScreen)