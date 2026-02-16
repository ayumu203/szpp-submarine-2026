package domain

type GameRepository interface {
	Save(game *Game) error
	FindByID(gameID string) (*Game, bool, error)
}

type TurnLogRepository interface {
	Append(gameID string, log TurnLog) error
	FindByGameID(gameID string) ([]TurnLog, error)
}

type PredictionRepository interface {
	Save(gameID, playerID string, board PredictionBoard) error
	Find(gameID, playerID string) (PredictionBoard, bool, error)
}

type PlayerGamesIndexRepository interface {
	AddGame(playerID, gameID string) error
	ListGames(playerID string) ([]string, error)
}
