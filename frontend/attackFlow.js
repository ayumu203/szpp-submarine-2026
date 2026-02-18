// 攻撃の処理を行うプログラム
const attackFlowState = {
    phase: "idle",  // idle | selectAttacker | selectTarget | submitting | opponentTurn
    selectedAttacker: null,
    selectedTarget: null,  // {x, y}
    candidateTargets: [],  // [{x,y}, ...]
    hasAttackedThisTurn: false
};

const attackTurnState = {
    viewerPlayerId: "hogehoge1",
    currentPlayerId: null
}

let attackAllySubmarines = [];  // [{x,y,hp,sunk,id}, ...]
let currentGameId = null;

function isMyAttackTurn() {
    return attackTurnState.currentPlayerId === attackTurnState.viewerPlayerId;
}

function attackCellToPosition($cell) {
    const row = $cell.parent().index();
    const col = $cell.index();
    if (row < 1 || row > BOARD_SIZE || col < 1 || col > BOARD_SIZE) return null;
    return { x: col, y: row };
}

function attackPositionToCell(pos) {
    return $("#field tr").eq(pos.y).children("td").eq(pos.x);
}

function attackContainsPos(list, pos) {
    return list.some(p => p.x === pos.x && p.y === pos.y);
}

function attackInBoard(x, y) {
    return x >= 1 && x <= BOARD_SIZE && y >= 1 && y <= BOARD_SIZE;
}

function attackIsAllySubmarineAt(pos) {
    return attackAllySubmarines.some(s => !s.sunk && s.x === pos.x && s.y === pos.y);
}

function getNeighbors8(pos) {
    const result = [];
    for (let dy = -1; dy <= 1; dy++) {
        for (let dx = -1; dx <= 1; dx++) {
            if (dx === 0 && dy === 0) continue;
            const nx = pos.x + dx;
            const ny = pos.y + dy;
            if (attackInBoard(nx, ny)) {
                result.push({ x: nx, y: ny });
            }
        }
    }
    return result;
}

// 攻撃先の候補を計算する（選択した1隻の周囲8マスのみ）
function getAttackableCells(selectedAttacker) {
    if (!selectedAttacker) {
        return [];
    }

    return getNeighbors8(selectedAttacker).filter(pos => !attackIsAllySubmarineAt(pos));
}

/* UI設計 */
function clearAttackHighlights() {
    $("#field td").removeClass("attack-clickable attack-selected-target");
}

function renderAttackHighlights() {
    clearAttackHighlights();

    if (attackFlowState.phase === "selectAttacker") {
        attackAllySubmarines
            .filter(s => !s.sunk)
            .forEach(s => {
                attackPositionToCell({ x: s.x, y: s.y }).addClass("attack-clickable");
            });

        if (attackFlowState.selectedAttacker) {
            attackPositionToCell(attackFlowState.selectedAttacker).addClass("attack-selected-target");
        }
        return;
    }

    if (attackFlowState.phase !== "selectTarget") return;

    attackFlowState.candidateTargets.forEach(p => {
        attackPositionToCell(p).addClass("attack-clickable");
    });

    if (attackFlowState.selectedTarget) {
        attackPositionToCell(attackFlowState.selectedTarget).addClass("attack-selected-target");
    }
}

function attackSetButtonEnabled(selector, enabled) {
    $(selector).prop("disabled", !enabled);
}

function updateAttackClickableControlsByPhase() {
    // 一旦全部無効化
    ["#btn-attack", "#btn-move", "#btn-back", "#btn-display", "#btn-apply"].forEach(id => {
        attackSetButtonEnabled(id, false);
    });

    if (attackFlowState.phase === "selectAttacker" || attackFlowState.phase === "selectTarget") {
        attackSetButtonEnabled("#btn-back", true);
        attackSetButtonEnabled("#btn-display", true);
        attackSetButtonEnabled("#btn-apply", true);
        return;
    }

    if (attackFlowState.phase === "idle") {
        attackSetButtonEnabled("#btn-attack", true);
        attackSetButtonEnabled("#btn-move", true);
        attackSetButtonEnabled("#btn-display", true);
        return;
    }

    if (attackFlowState.phase === "opponentTurn") {
        attackSetButtonEnabled("#btn-display", true);
    }
}

// 実際の攻撃処理
function startAttackFlow() {
    if (!isMyAttackTurn()) return;
    if (attackFlowState.hasAttackedThisTurn) return;

    attackFlowState.phase = "selectAttacker";
    attackFlowState.selectedAttacker = null;
    attackFlowState.selectedTarget = null;
    attackFlowState.candidateTargets = [];
    renderAttackHighlights();
    updateAttackClickableControlsByPhase();
}

function cancelAttackFlow() {
    attackFlowState.phase = "idle";
    attackFlowState.selectedAttacker = null;
    attackFlowState.selectedTarget = null;
    attackFlowState.candidateTargets = [];
    updateAttackClickableControlsByPhase();
    clearAttackHighlights();
}

function handleFieldClickForAttack($cell) {
    const pos = attackCellToPosition($cell);
    if (!pos) return;

    if (attackFlowState.phase === "selectAttacker") {
        const attacker = attackAllySubmarines.find(
            s => !s.sunk && s.x === pos.x && s.y === pos.y
        );
        if (!attacker) return;

        attackFlowState.selectedAttacker = { id: attacker.id, x: attacker.x, y: attacker.y };
        renderAttackHighlights();
        return;
    }

    if (attackFlowState.phase === "selectTarget") {
        if (!attackContainsPos(attackFlowState.candidateTargets, pos)) return;
        attackFlowState.selectedTarget = pos;
        renderAttackHighlights();
    }
}

function confirmAttackStep1() {
    if (attackFlowState.phase !== "selectAttacker") return;
    if (!attackFlowState.selectedAttacker) return;

    attackFlowState.phase = "selectTarget";
    attackFlowState.candidateTargets = getAttackableCells(attackFlowState.selectedAttacker);
    attackFlowState.selectedTarget = null;

    renderAttackHighlights();
    updateAttackClickableControlsByPhase();
}

async function confirmAttackStep2() {
    if (attackFlowState.phase !== "selectTarget") return;
    if (!attackFlowState.selectedAttacker || !attackFlowState.selectedTarget) return;

    // APIに送信するリクエスト
    const payload = {
        gameId: currentGameId,
        playerId: attackTurnState.viewerPlayerId,
        actionType: "attack",
        // 必要なら attackerId も送る（API仕様次第）
        target: attackFlowState.selectedTarget
    };

    attackFlowState.phase = "submitting";
    updateAttackClickableControlsByPhase();

    try {
        if (typeof executeAction === "function") {
            await executeAction(payload);
        } else {
            console.log("Attack payload:", payload);
        }

        attackFlowState.hasAttackedThisTurn = true;
        attackFlowState.phase = "opponentTurn";
        attackFlowState.selectedAttacker = null;
        attackFlowState.selectedTarget = null;
        attackFlowState.candidateTargets = [];
        clearAttackHighlights();
    } catch (e) {
        attackFlowState.phase = "selectTarget";
        renderAttackHighlights();
        updateAttackClickableControlsByPhase();
        console.error("攻撃フローの送信に失敗しました: ", e);
    }
}

function bindAttackFlowEvents() {
    $("#btn-attack").off("click.attackFlow").on("click.attackFlow", function () {
        startAttackFlow();
    });

    $("#btn-back").off("click.attackFlow").on("click.attackFlow", function () {
        if (attackFlowState.phase === "selectTarget") {
            attackFlowState.phase = "selectAttacker";
            attackFlowState.selectedTarget = null;
            attackFlowState.candidateTargets = [];
            renderAttackHighlights();
            updateAttackClickableControlsByPhase();
            return;
        }

        if (attackFlowState.phase === "selectAttacker") {
            cancelAttackFlow();
        }
    });

    $("#btn-apply").off("click.attackFlow").on("click.attackFlow", async function () {
        if (attackFlowState.phase === "selectAttacker") {
            confirmAttackStep1();
            return;
        }
        if (attackFlowState.phase === "selectTarget") {
            await confirmAttackStep2();
        }
    });

    $("#field td").off("click.attackFlow").on("click.attackFlow", function () {
        handleFieldClickForAttack($(this));
    });
}

function syncAttackContextFromState(gameState) {
    currentGameId = gameState.gameId;
    attackTurnState.currentPlayerId = gameState.currentPlayerId;

    attackAllySubmarines = Object.entries(gameState.allyBoard.submarines).map(([id, submarine]) => ({
        id,
        x: submarine.x,
        y: submarine.y,
        hp: submarine.hp,
        sunk: submarine.sunk
    }));

    if (isMyAttackTurn()) {
        attackFlowState.hasAttackedThisTurn = false;
        if (attackFlowState.phase === "opponentTurn") {
            attackFlowState.phase = "idle";
        }
    } else {
        attackFlowState.phase = "opponentTurn";
    }

    updateAttackClickableControlsByPhase();
}