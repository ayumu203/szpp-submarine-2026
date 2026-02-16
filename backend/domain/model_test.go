package domain

import "testing"

func TestAttackWithinRange(t *testing.T) {
	board := NewBoard()
	board.PlaceSubmarine("p1", "p1-sub-1", Position{X: 2, Y: 2})
	board.PlaceSubmarine("p2", "p2-sub-1", Position{X: 3, Y: 3})

	report, errCode := board.Attack("p1", Position{X: 3, Y: 3})
	if errCode != "" {
		t.Fatalf("unexpected error code: %s", errCode)
	}
	if report != ReportHit {
		t.Fatalf("expected report %s, got %s", ReportHit, report)
	}
}

func TestPlayerMoveWithinRange(t *testing.T) {
	tests := []struct {
		name      string
		distance  int
		expectErr ErrorCode
		expect    MoveReportType
	}{
		{name: "move 1 cell", distance: 1, expect: MoveReportSuccess},
		{name: "move 2 cells", distance: 2, expect: MoveReportSuccess},
		{name: "invalid distance", distance: 3, expectErr: ErrorInvalidMoveDistance},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := NewBoard()
			board.PlaceSubmarine("p1", "p1-sub-1", Position{X: 2, Y: 2})
			report, errCode := board.MoveSubmarine("p1", DirectionEast, tt.distance)
			if tt.expectErr != "" {
				if errCode != tt.expectErr {
					t.Fatalf("expected err %s, got %s", tt.expectErr, errCode)
				}
				return
			}
			if errCode != "" {
				t.Fatalf("unexpected error code: %s", errCode)
			}
			if report != tt.expect {
				t.Fatalf("expected report %s, got %s", tt.expect, report)
			}
		})
	}
}

func TestSunkStateHandling(t *testing.T) {
	board := NewBoard()
	board.PlaceSubmarine("p1", "p1-sub-1", Position{X: 2, Y: 2})
	board.PlaceSubmarine("p2", "p2-sub-1", Position{X: 3, Y: 3})

	for i := 0; i < 3; i++ {
		report, errCode := board.Attack("p1", Position{X: 3, Y: 3})
		if errCode != "" {
			t.Fatalf("unexpected error code: %s", errCode)
		}
		if i < 2 && report != ReportHit {
			t.Fatalf("expected hit before sunk, got %s", report)
		}
		if i == 2 && report != ReportHitAndSunk {
			t.Fatalf("expected hitAndSunk on final hit, got %s", report)
		}
	}

	if !board.Submarines["p2-sub-1"].IsSunk() {
		t.Fatal("expected submarine to be sunk")
	}
}

func TestAttackEnemyCellReturnsHitEvenWhenOutOfRange(t *testing.T) {
	board := NewBoard()
	board.PlaceSubmarine("p1", "p1-sub-1", Position{X: 1, Y: 1})
	board.PlaceSubmarine("p2", "p2-sub-1", Position{X: 5, Y: 5})

	report, errCode := board.Attack("p1", Position{X: 5, Y: 5})
	if errCode != "" {
		t.Fatalf("unexpected error code: %s", errCode)
	}
	if report != ReportHit {
		t.Fatalf("expected report %s, got %s", ReportHit, report)
	}
}
