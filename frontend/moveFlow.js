/**
 * 実装方針
 * 
 * 移動専用の状態管理 moveFlowState を導入
 * 
 * フェーズを明確化
 * 
 * 1. idle
 * 2. selectSource（移動する潜水艦の選択）
 * 3. selectDestination（移動先の選択）
 * 4. submitting
 * 5. opponentTurn
 * 
 * 決定 ボタンを2回使う
 * 1回目: 潜水艦選択の確定
 * 2回目: 移動先選択の確定（この時点でAPI送信）
 * クリック制限
 * selectSource: 潜水艦セル + 戻る + 予測表示切替 + 決定 だけ有効
 * selectDestination: 移動可能セル + 戻る + 予測表示切替 + 決定 だけ有効
 * 1ターンに移動できるのは1隻のみ
 * 移動成功後 hasMovedThisTurn = true
 * currentPlayerId !== viewerPlayerId を検知したら opponentTurn
 * 自分のターンに戻ったら hasMovedThisTurn = false にリセット
 */

const moveFlowState = {
    phase: "idle",  // idle | selectSource | selectDestination | submitting | opponentTurn
    selectedSource: null,  // {x, y}
    selectedDestination: null,  // {x, y}
    candidateDestinations: [],  // [{x, y}, ...]
    hasMovedThisTurn: false
}

// 画面側の現在ターン情報（state API反映時に更新）
const turnState = {
    viewerPlayerId: "playerA",
    currentPlayerId: null
}

// ally潜水艦の位置情報の初期化（state API反映時に更新）
let allySubmarines = [];

/**
 * @description 現在が自分のターンか
 * @returns {boolean} trueなら自分のターン
 */
function isMyTurn() {
    return turnState.currentPlayerId === turnState.viewerPlayerId;
}

/**
 * @description 選んだ2つのセルが同じか
 * @param {object} a 選んだセル1
 * @param {object} b 選んだセル2
 * @returns {boolean} trueなら同じ
 */
function isSameCell(a, b) {
    return a && b && a.x === b.x && a.y === b.y;
}

/**
 * @description 選んだセルの座標情報を取得する
 * @returns {{x: number, y: number}} {x座標, y座標}
 */
function cellToPosition($cell) {
    const row = $cell.parent().index();
    const col = $cell.index();
    if (row < 1 || row > BOARD_SIZE || col < 1 || col > BOARD_SIZE) return null;
    return {x: col, y: row};
}

/**
 * @description 座標からセルを特定する
 * @param {{number, number}} pos 座標情報
 */
function positionToCell(pos) {
    return $("#field tr").eq(pos.y).children("td").eq(pos.x);
}

function containsPos(list, pos) {
    return list.some(p => p.x === pos.x && p.y === pos.y);
}

/**
 * @description その座標に味方潜水艦がいるか
 * @param {{x: number, y: number}} pos 座標情報
 */
function isAllySubmarineAt(pos) {
    return allySubmarines.some(s => !s.sunk && s.x === pos.x && s.y === pos.y);
}

/**
 * @description 移動先の潜水艦の移動先候補を計算する
 * ルール: 東西南北1～2マス、盤外NG、占有セルNG、経路上に撃沈艦があればその先不可
 * @param {object} source 選択した潜水艦
 * @param {object} submarines 敵味方すべての潜水艦
 * @returns {[{x: number, y: number}, ...]} 移動可能な座標のリスト
 */
function calculateMoveCandidates(source, submarines) {
    // 方角を座標に変換
    const dirs = [
        {dx: 0, dy: -1},  // 北
        {dx: 0, dy: 1},   // 南
        {dx: 1, dy: 0},    // 東
        {dx: -1, dy: 0}   // 西
    ];

    // すでに潜水艦がいるセルの集合
    const occupied = new Set(submarines.map(s => `${s.x},${s.y}`));

    // 撃沈した潜水艦がいるセルの集合
    const sunk = new Set(submarines.filter(s => s.sunk).map(s => `${s.x},${s.y}`));

    const result = [];

    // すべての方角に対してチェック
    for (const d of dirs) {
        // 2マスまでなら動ける
        for (let step = 1; step <= 2; step++) {
            // 変数nx, nyを宣言し、方角dにstepだけ進んだ先の座標を格納する
            const nx =  source.x + d.dx * step;
            const ny =  source.y + d.dy * step;
            
            if (nx < 1 || nx > BOARD_SIZE || ny < 1 || ny > BOARD_SIZE) break;

            const key = `${nx},${ny}`;

            // 経路上に撃沈艦があれば進行停止
            if (sunk.has(key)) break;

            // 目的地が占有セルなら不可
            if (occupied.has(key)) break;

            result.push({x: nx, y: ny});
        }
    }
    return result;
}

/**
 * @description 移動先がもとの場所から「どの方角に」「何マス離れているか」を計算する
 * @param {{x: number, y: number}} from 元々いた座標
 * @param {{x: number, y: number}} to 移動先の座標
 * @returns {{direction: string, distance: number}} {方角, 移動距離}
 */
function toDirectionAndDistance(from, to) {
    // dx, dyにはそれぞれ「移動先の座標」-「元々いた座標」を代入する
    const dx = to.x - from.x; // 頑張って書く
    const dy = to.y - from.y; // 頑張って書く

    if (dx !== 0 && dy !== 0) {
        throw new Error("斜め移動はできません");
    }

    const distance = Math.abs(dx) + Math.abs(dy);

    if (distance < 1 || distance > 2) {
        throw new Error("移動できるのは1マスまたは2マスのみです");
    }

    // directionには「どの方角に移動したのか」を文字列で格納する
    let direction = "";
    if (dx > 0) {
      direction = "east";
    } else if (dx < 0) {
      direction = "west";
    } else if (dy > 0) {
      direction = "north";
    } else if (dy < 0) {
      direction = "south";
    }

    return {direction, distance}
}

function clearMoveHighlights() {
  $("#field td")
    .removeClass("move-clickable move-source move-destination move-candidate move-disabled");
}

function renderMoveHighlights() {
  clearMoveHighlights();

  if (moveFlowState.phase === "selectSource") {
    allySubmarines
      .filter(s => !s.sunk)
      .forEach(s => positionToCell({ x: s.x, y: s.y }).addClass("move-clickable"));

    if (moveFlowState.selectedSource) {
      positionToCell(moveFlowState.selectedSource).addClass("move-source");
    }
    return;
  }
  if (moveFlowState.phase === "selectDestination") {
    if (moveFlowState.selectedSource) {
      positionToCell(moveFlowState.selectedSource).addClass("move-source");
    }

    moveFlowState.candidateDestinations.forEach(p => {
      positionToCell(p).addClass("move-clickable move-candidate");
    });

    if (moveFlowState.selectedDestination) {
      positionToCell(moveFlowState.selectedDestination).addClass("move-destination");
    }
  }
}

function setButtonEnabled(selector, enabled) {
  $(selector).prop("disabled", !enabled);
}

function updateClickableControlsByPhase() {
  // まず全部無効
  ["#btn-attack", "#btn-move", "#btn-back", "#btn-display", "#btn-apply"].forEach(id => {
    setButtonEnabled(id, false);
  });

  if (moveFlowState.phase === "selectSource" || moveFlowState.phase === "selectDestination") {
    setButtonEnabled("#btn-back", true);
    setButtonEnabled("#btn-display", true);
    setButtonEnabled("#btn-apply", true);
    return;
  }

  if (moveFlowState.phase === "idle") {
    setButtonEnabled("#btn-attack", true);
    setButtonEnabled("#btn-move", true);
    setButtonEnabled("#btn-display", true);
    return;
  }

  if (moveFlowState.phase === "opponentTurn") {
    setButtonEnabled("#btn-display", true);
    return;
  }
}

function startMoveFlow() {
  if (!isMyTurn()) return;
  if (moveFlowState.hasMovedThisTurn) return; // 1ターン1隻

  moveFlowState.phase = "selectSource";
  moveFlowState.selectedSource = null;
  moveFlowState.selectedDestination = null;
  moveFlowState.candidateDestinations = [];

  updateClickableControlsByPhase();
  renderMoveHighlights();
}

function cancelMoveFlow() {
  moveFlowState.phase = "idle";
  moveFlowState.selectedSource = null;
  moveFlowState.selectedDestination = null;
  moveFlowState.candidateDestinations = [];
  updateClickableControlsByPhase();
  clearMoveHighlights();
}

function handleFieldClickForMove($cell) {
  const pos = cellToPosition($cell);
  if (!pos) return;

  if (moveFlowState.phase === "selectSource") {
    if (!isAllySubmarineAt(pos)) return;
    moveFlowState.selectedSource = pos;
    renderMoveHighlights();
    return;
  }

  if (moveFlowState.phase === "selectDestination") {
    if (!containsPos(moveFlowState.candidateDestinations, pos)) return;
    moveFlowState.selectedDestination = pos;
    renderMoveHighlights();
  }
}

function confirmMoveStep1() {
  // 潜水艦選択の確定
  if (!moveFlowState.selectedSource) return;

  moveFlowState.candidateDestinations = calculateMoveCandidates(
    moveFlowState.selectedSource,
    allySubmarines
  );
  moveFlowState.selectedDestination = null;
  moveFlowState.phase = "selectDestination";

  updateClickableControlsByPhase();
  renderMoveHighlights();
}

async function confirmMoveStep2() {
  // 移動先選択の確定 + 送信
  if (!moveFlowState.selectedSource || !moveFlowState.selectedDestination) return;

  moveFlowState.phase = "submitting";
  updateClickableControlsByPhase();

  const { direction, distance } = toDirectionAndDistance(
    moveFlowState.selectedSource,
    moveFlowState.selectedDestination
  );

  // [must] window.currentGameId 未設定時の防御的チェック
  if (typeof window === "undefined" || window.currentGameId == null) {
    console.error("[moveFlow] window.currentGameId が未設定のため、移動アクションを送信できません。");
    // フェーズと UI を「移動先選択」状態に戻す
    moveFlowState.phase = "selectDestination";
    updateClickableControlsByPhase();
    renderMoveHighlights();
    return;
  }
  const payload = {
    gameId: window.currentGameId,
    playerId: turnState.viewerPlayerId,
    actionType: "move",
    direction,
    distance
  };
  try {
    // ここは既存のAPI関数へ接続（例）
    // const result = await executeAction(payload);

    // 成功時の扱い
    moveFlowState.hasMovedThisTurn = true;
    moveFlowState.phase = "opponentTurn";
    moveFlowState.selectedSource = null;
    moveFlowState.selectedDestination = null;
    moveFlowState.candidateDestinations = [];
    clearMoveHighlights();

    // API結果 or 最新stateでターン更新（例）
    // turnState.currentPlayerId = result.nextPlayerId;
    // await refreshStateAndBoard();
  } catch (e) {
    // 失敗時は移動先選択フェーズへ戻す
    moveFlowState.phase = "selectDestination";
    updateClickableControlsByPhase();
    renderMoveHighlights();
  }
}

function bindMoveFlowEvents() {
  $("#btn-move").off("click.moveFlow").on("click.moveFlow", startMoveFlow);

  $("#btn-back").off("click.moveFlow").on("click.moveFlow", function () {
    if (moveFlowState.phase === "selectDestination") {
      moveFlowState.phase = "selectSource";
      moveFlowState.selectedDestination = null;
      moveFlowState.candidateDestinations = [];
      updateClickableControlsByPhase();
      renderMoveHighlights();
      return;
    }
    if (moveFlowState.phase === "selectSource") {
      cancelMoveFlow();
    }
  });

  $("#btn-apply").off("click.moveFlow").on("click.moveFlow", async function () {
    if (moveFlowState.phase === "selectSource") {
      confirmMoveStep1();
      return;
    }
    if (moveFlowState.phase === "selectDestination") {
      await confirmMoveStep2();
    }
  });

  $("#field td").off("click.moveFlow").on("click.moveFlow", function () {
    handleFieldClickForMove($(this));
  });
}

function syncMoveContextFromState(gameState, viewerPlayerId) {
  turnState.viewerPlayerId = viewerPlayerId;
  turnState.currentPlayerId = gameState.currentPlayerId;

  allySubmarines = Object.entries(gameState.allyBoard.submarines).map(([id, submarine]) => ({
    id,
    x: submarine.x,
    y: submarine.y,
    hp: submarine.hp,
    sunk: submarine.sunk
  }));

  if (isMyTurn()) {
    moveFlowState.hasMovedThisTurn = false;
    if (moveFlowState.phase === "opponentTurn") {
      moveFlowState.phase = "idle";
    }
  } else {
    moveFlowState.phase = "opponentTurn";
  }

  updateClickableControlsByPhase();
}