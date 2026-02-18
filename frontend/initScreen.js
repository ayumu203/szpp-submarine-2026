// ゲーム開始時に、画面の初期化を行う処理を記述する
const uiState = {
    mode: "idle",  // idle | attack | moveの3種類
    displayMode: "ally", // ally | prediction
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
 * @description 現在の手番ステータスをヘッダーに表示する
 * @param {object} gameState ゲーム状態
 * @param {string} viewerPlayerId 表示中プレイヤーID
 */
function renderTurnStatus(gameState, viewerPlayerId) {
    const statusElement = document.getElementById("turn-status");
    if (!statusElement) {
        return;
    }

    const isMyTurn = gameState.currentPlayerId === viewerPlayerId;
    statusElement.textContent = isMyTurn
        ? "あなたのターンです"
        : "相手ターンです";
}

/**
 * @description ゲーム開始時の画面を作成する
 */
async function initializeScreen() {
    // テーブルが存在しなかったらreturn
    if (!validateRequiredElements()) {
        return;
    }

    // UIをリセットする関数を呼び出す
    resetUiState();
    // main.jsの関数を使って、潜水艦を表示する
    renderSubmarines();

    bindDisplayToggle();
    await renderDisplayMode();

    const data = await getMock();
    const gameState = data.State.GetGameStateResponse;
    const initializeResponse = data.Initialize?.InitializeGameResponse;
    const viewerPlayerId = data.State.GetGameStateRequest.viewerPlayerId;
    const synchronizedGameState = {
        ...gameState,
        currentPlayerId: initializeResponse?.currentPlayerId ?? gameState.currentPlayerId
    };

    renderTurnStatus(synchronizedGameState, viewerPlayerId);

    if (typeof syncAttackContextFromState === "function") {
        attackTurnState.viewerPlayerId = viewerPlayerId;
        syncAttackContextFromState(synchronizedGameState);
    }
    if (typeof bindAttackFlowEvents === "function") {
        bindAttackFlowEvents();
    }

    if (typeof syncMoveContextFromState === "function") {
        syncMoveContextFromState(synchronizedGameState, viewerPlayerId);
    }
    if (typeof bindMoveFlowEvents === "function") {
        bindMoveFlowEvents();
    }
}

$(initializeScreen)