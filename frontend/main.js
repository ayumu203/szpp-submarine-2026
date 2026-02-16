$(function () {
    const state = {
        apiBaseUrl: "http://localhost:8081",
        gameId: "",
        viewerPlayerId: "human1",
    };

    function normalizeBaseUrl(url) {
        return url.replace(/\/$/, "");
    }

    function yLabelToNumber(label) {
        const labels = ["A", "B", "C", "D", "E"];
        return labels.indexOf((label || "").toUpperCase()) + 1;
    }

    function numberToYLabel(number) {
        const labels = ["A", "B", "C", "D", "E"];
        return labels[number - 1] || String(number);
    }

    function setMessage(selector, message, isError) {
        const element = $(selector);
        element.text(message || "");
        element.toggleClass("error", Boolean(isError));
    }

    function getConfigFromInputs() {
        state.apiBaseUrl = normalizeBaseUrl($("#apiBaseUrl").val().trim());
        state.viewerPlayerId = $("#actionPlayerId").val().trim() || $("#playerAId").val().trim();
    }

    async function callApi(path, method, body) {
        const options = {
            url: `${state.apiBaseUrl}${path}`,
            method,
            contentType: "application/json",
        };
        if (body) {
            options.data = JSON.stringify(body);
        }
        return $.ajax(options);
    }

    function createBoardTable(cells, submarines) {
        const table = $("<table>");
        const header = $("<tr>").append("<th></th>");
        for (let x = 1; x <= 5; x += 1) {
            header.append(`<th>${x}</th>`);
        }
        table.append(header);

        const labels = ["A", "B", "C", "D", "E"];
        for (let y = 0; y < 5; y += 1) {
            const row = $("<tr>").append(`<th>${labels[y]}</th>`);
            for (let x = 0; x < 5; x += 1) {
                const cell = $("<td>");
                const submarineId = cells?.[y]?.[x];
                if (submarineId && submarines[submarineId]) {
                    const submarine = submarines[submarineId];
                    const ownerClass = submarine.ownerId === state.viewerPlayerId ? "ally" : "enemy";
                    const sunkClass = submarine.sunk ? "sunk" : "";
                    cell.addClass(`${ownerClass} ${sunkClass}`);
                    cell.html(`<div class="cellMain">${submarine.ownerId}</div><div class="cellSub">HP:${submarine.hp}</div>`);
                }
                row.append(cell);
            }
            table.append(row);
        }

        return table;
    }

    function renderLogs(logs) {
        const list = $("#logList");
        list.empty();
        if (!logs || logs.length === 0) {
            list.append("<li>ログはまだありません。</li>");
            return;
        }
        const sorted = [...logs].sort((a, b) => b.turn - a.turn);
        sorted.forEach((log) => {
            const parts = [`T${log.turn}`, log.playerId, log.actionType];
            if (log.target) {
                parts.push(`target(${log.target.x},${numberToYLabel(log.target.y)})`);
            }
            if (log.direction) {
                parts.push(`dir:${log.direction}`);
            }
            if (log.distance) {
                parts.push(`dist:${log.distance}`);
            }
            if (log.attackReport) {
                parts.push(`attack:${log.attackReport}`);
            }
            if (log.moveReport) {
                parts.push(`move:${log.moveReport}`);
            }
            if (log.errorCode) {
                parts.push(`error:${log.errorCode}`);
            }
            list.append(`<li>${parts.join(" | ")}</li>`);
        });
    }

    function renderState(data) {
        $("#gameId").text(data.gameId || "-");
        $("#turn").text(String(data.turn ?? "-"));
        $("#status").text(data.status || "-");
        $("#currentPlayer").text(data.currentPlayerId || "-");

        const boardElement = $("#board");
        boardElement.empty();
        boardElement.append(createBoardTable(data.board?.cells, data.board?.submarines || {}));

        renderLogs(data.logs || []);
    }

    async function fetchState() {
        if (!state.gameId) {
            setMessage("#configMessage", "先にゲームを初期化してください。", true);
            return;
        }
        getConfigFromInputs();
        try {
            const data = await callApi(`/state?gameId=${encodeURIComponent(state.gameId)}&viewerPlayerId=${encodeURIComponent(state.viewerPlayerId)}`, "GET");
            renderState(data);
            setMessage("#configMessage", "状態更新に成功しました。", false);
        } catch (error) {
            const message = error.responseText || error.statusText || "状態取得に失敗しました。";
            setMessage("#configMessage", message, true);
        }
    }

    $("#configForm").on("submit", async function (event) {
        event.preventDefault();
        getConfigFromInputs();

        const requestBody = {
            playerAId: $("#playerAId").val().trim(),
            playerBId: $("#playerBId").val().trim(),
        };

        try {
            const data = await callApi("/initialize", "POST", requestBody);
            state.gameId = data.gameId;
            $("#gameId").text(state.gameId);
            setMessage("#configMessage", `ゲームを作成しました: ${state.gameId}`, false);
            await fetchState();
        } catch (error) {
            const message = error.responseText || error.statusText || "初期化に失敗しました。";
            setMessage("#configMessage", message, true);
        }
    });

    $("#actionType").on("change", function () {
        const actionType = $(this).val();
        const isAttack = actionType === "attack";
        $(".attackField").toggleClass("hidden", !isAttack);
        $(".moveField").toggleClass("hidden", isAttack);
    });

    $("#actionForm").on("submit", async function (event) {
        event.preventDefault();
        getConfigFromInputs();

        if (!state.gameId) {
            setMessage("#actionMessage", "先にゲームを初期化してください。", true);
            return;
        }

        const actionType = $("#actionType").val();
        const requestBody = {
            gameId: state.gameId,
            playerId: $("#actionPlayerId").val().trim(),
            actionType,
        };

        if (actionType === "attack") {
            requestBody.target = {
                x: Number($("#targetX").val()),
                y: yLabelToNumber($("#targetY").val()),
            };
        } else {
            requestBody.direction = $("#direction").val();
            requestBody.distance = Number($("#distance").val());
        }

        try {
            const data = await callApi("/action", "POST", requestBody);
            const report = data.errorCode || data.attackReport || data.moveReport || "ok";
            setMessage("#actionMessage", `行動結果: ${report}`, Boolean(data.errorCode));
            await fetchState();
        } catch (error) {
            const message = error.responseText || error.statusText || "行動実行に失敗しました。";
            setMessage("#actionMessage", message, true);
        }
    });

    $("#refreshButton").on("click", fetchState);
    $("#actionType").trigger("change");
});
