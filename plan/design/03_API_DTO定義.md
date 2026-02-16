# API DTO定義

## 方針
- ドメインモデルとは分離し、HTTP入出力専用のDTOを定義する。
- `playerId` 命名に統一する。
- 座標系は `x,y` ともに `1..5` を採用する。
- `POST /action` は1リクエスト内で「人間手 + 必要時CPU手」まで処理して応答する。
- `GetGameStateResponse.board` は相手情報を含めて全公開（マスクなし）とする。

## Initialize
### Request: `InitializeGameRequest`
- `playerAId: string`
- `playerBId: string`

### Response: `InitializeGameResponse`
- `gameId: string`
- `status: string`
- `turn: number`
- `currentPlayerId: string`

## Action
### Request: `ExecuteActionRequest`
- `gameId: string`
- `playerId: string`
- `actionType: "attack" | "move"`
- `target?: { x: number, y: number }`
- `direction?: "north" | "south" | "east" | "west"`
- `distance?: number` (`1` or `2`)

### Response: `ExecuteActionResponse`
- `gameId: string`
- `turn: number`
- `attackReport?: "miss" | "hit" | "hitAndSunk" | "waveHigh"`
- `moveReport?: "moveSuccess" | "moveBlocked"`
- `errorCode?: "invalidTurn" | "invalidAction" | "invalidTarget" | "invalidMoveDistance" | "outOfBoard"`
- `nextPlayerId: string`
- `winnerId?: string`
- `status: "waiting" | "inProgress" | "finished"`

## State
### Request: `GetGameStateRequest`
- `gameId: string`
- `viewerPlayerId: string`

### Response: `GetGameStateResponse`
- `gameId: string`
- `turn: number`
- `status: string`
- `currentPlayerId: string`
- `opponentId: string`
- `board: BoardViewDto`
- `predictionBoard: PredictionBoardDto`
- `logs: TurnLogDto[]`

## 子DTO
### `BoardViewDto`
- `cells: (string | null)[][]` (5x5)
- `submarines: Record<string, { ownerId: string, x: number, y: number, hp: number, sunk: boolean }>`

### `PredictionBoardDto`
- `scoreGrid: number[][]` (5x5)
- `possibleEnemyCount: number[][]` (5x5)
- `updatedAt: string`

### `TurnLogDto`
- `turn: number`
- `playerId: string`
- `actionType: "attack" | "move"`
- `target?: { x: number, y: number }`
- `direction?: "north" | "south" | "east" | "west"`
- `distance?: number`
- `attackReport?: "miss" | "hit" | "hitAndSunk" | "waveHigh"`
- `moveReport?: "moveSuccess" | "moveBlocked"`
- `errorCode?: "invalidTurn" | "invalidAction" | "invalidTarget" | "invalidMoveDistance" | "outOfBoard"`
- `createdAt: string`
