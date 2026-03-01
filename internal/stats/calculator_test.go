package stats

import (
	"testing"

	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/parser"
)

// Boundary Value Testing (BVT) for stats/calculator.go
// Tests focus on critical decision boundaries and equivalence class boundaries

func TestCalculatorNewCalculator(t *testing.T) {
	calc := NewCalculator()
	if calc == nil {
		t.Fatal("NewCalculator returned nil")
	}
}

func TestCalculatorCalculateEmptyHands(t *testing.T) {
	calc := NewCalculator()
	stats := calc.Calculate([]*parser.Hand{}, 0)

	if stats == nil {
		t.Fatal("Calculate returned nil for empty hands")
	}
	if stats.TotalHands != 0 {
		t.Errorf("expected 0 total hands, got %d", stats.TotalHands)
	}
}

func TestCalculatorInvestedAmount(t *testing.T) {
	calc := NewCalculator()
	hand := &parser.Hand{
		Players: make(map[int]*parser.PlayerHandInfo),
	}

	// Boundary: Player not in hand (missing)
	amount := calc.investedAmount(hand, 0)
	if amount != 0 {
		t.Errorf("boundary - missing player: expected 0, got %d", amount)
	}

	// Boundary: Player with no actions (empty)
	hand.Players[0] = &parser.PlayerHandInfo{
		Actions: []parser.PlayerAction{},
	}
	amount = calc.investedAmount(hand, 0)
	if amount != 0 {
		t.Errorf("boundary - empty actions: expected 0, got %d", amount)
	}

	// Boundary: Player with exactly 1 action
	hand.Players[0].Actions = []parser.PlayerAction{
		{Amount: 20},
	}
	amount = calc.investedAmount(hand, 0)
	if amount != 20 {
		t.Errorf("boundary - 1 action: expected 20, got %d", amount)
	}

	// Normal case: Multiple actions
	hand.Players[0].Actions = []parser.PlayerAction{
		{Amount: 10},
		{Amount: 20},
		{Amount: 30},
	}
	amount = calc.investedAmount(hand, 0)
	if amount != 60 {
		t.Errorf("normal case: expected 60, got %d", amount)
	}
}

func TestEnsurePositionStats(t *testing.T) {
	calc := NewCalculator()
	stats := &Stats{
		ByPosition: make(map[parser.Position]*PositionStats),
	}

	// Boundary: First call creates new
	pos1 := calc.ensurePositionStats(stats, parser.PosBTN)
	if pos1 == nil {
		t.Fatal("ensurePositionStats returned nil")
	}

	// Boundary: Second call retrieves same (idempotent)
	pos2 := calc.ensurePositionStats(stats, parser.PosBTN)
	if pos1 != pos2 {
		t.Error("boundary - idempotency: should return same instance")
	}

	// Different position creates new
	pos3 := calc.ensurePositionStats(stats, parser.PosSB)
	if pos1 == pos3 {
		t.Error("boundary - different positions: should have different instances")
	}
}

func TestBBAmountFromHand(t *testing.T) {
	// Boundary: Invalid BB seat (-1)
	hand := &parser.Hand{
		BBSeat:  -1,
		Players: make(map[int]*parser.PlayerHandInfo),
	}
	bb := bbAmountFromHand(hand)
	if bb != 0 {
		t.Errorf("boundary - invalid BB seat: expected 0, got %d", bb)
	}

	// Boundary: Valid BB seat but no player
	hand.BBSeat = 1
	bb = bbAmountFromHand(hand)
	if bb != 0 {
		t.Errorf("boundary - no BB player: expected 0, got %d", bb)
	}

	// Boundary: Valid BB seat and player but no BB action
	hand.Players[1] = &parser.PlayerHandInfo{
		Actions: []parser.PlayerAction{
			{Amount: 20, Action: parser.ActionCall},
		},
	}
	bb = bbAmountFromHand(hand)
	if bb != 0 {
		t.Errorf("boundary - no BB action: expected 0, got %d", bb)
	}

	// Boundary: Valid BB with exactly 1 BB action
	hand.Players[1].Actions = []parser.PlayerAction{
		{Amount: 20, Action: parser.ActionBlindBB, Street: parser.StreetPreFlop},
	}
	bb = bbAmountFromHand(hand)
	if bb != 20 {
		t.Errorf("boundary - 1 BB action: expected 20, got %d", bb)
	}

	// Normal case: Multiple actions with BB
	hand.Players[1].Actions = []parser.PlayerAction{
		{Amount: 10, Action: parser.ActionCall},
		{Amount: 20, Action: parser.ActionBlindBB},
		{Amount: 30, Action: parser.ActionRaise},
	}
	bb = bbAmountFromHand(hand)
	if bb != 20 {
		t.Errorf("normal case: expected 20, got %d", bb)
	}
}

// BOUNDARY VALUE TEST SUITE for bucketByBBMultiple
// Critical boundaries: 0, 2.5x, 4.0x, 6.0x, 10.0x
func TestBucketByBBMultiple(t *testing.T) {
	tests := []struct {
		name     string
		amount   int
		bb       int
		expected RangeActionBucket
		note     string
	}{
		// Lower boundary (0)
		{"0 (lower)", 0, 20, RangeActionCheck, "min"},

		// Boundary at 2.5x: amount/bb <= 2.5 → BetSmall
		{"exactly 2.5x", 50, 20, RangeActionBetSmall, "at-boundary"},
		{"2.4x (49/20)", 49, 20, RangeActionBetSmall, "just-below"},
		{"2.6x (52/20)", 52, 20, RangeActionBetHalf, "just-above"}, // 2.6 > 2.5, enters BetHalf

		// Boundary at 4.0x: amount/bb <= 4.0 → BetHalf
		{"exactly 4.0x", 80, 20, RangeActionBetHalf, "at-boundary"},
		{"3.9x (78/20)", 78, 20, RangeActionBetHalf, "just-below"},
		{"4.1x (82/20)", 82, 20, RangeActionBetTwoThird, "just-above"}, // 4.1 > 4.0

		// Boundary at 6.0x: amount/bb <= 6.0 → BetTwoThird
		{"exactly 6.0x", 120, 20, RangeActionBetTwoThird, "at-boundary"},
		{"5.9x (118/20)", 118, 20, RangeActionBetTwoThird, "just-below"},
		{"6.1x (122/20)", 122, 20, RangeActionBetPot, "just-above"}, // 6.1 > 6.0

		// Boundary at 10.0x: amount/bb <= 10.0 → BetPot
		{"exactly 10.0x", 200, 20, RangeActionBetPot, "at-boundary"},
		{"9.9x (198/20)", 198, 20, RangeActionBetPot, "just-below"},
		{"10.1x (202/20)", 202, 20, RangeActionBetOver, "just-above"}, // 10.1 > 10.0

		// Upper boundary (very large)
		{"100x (extreme)", 2000, 20, RangeActionBetOver, "max"},

		// Edge case: zero BB defaults to 20
		{"zero BB edge case", 50, 0, RangeActionBetSmall, "default-bb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bucketByBBMultiple(tt.amount, tt.bb)
			if result != tt.expected {
				actual := float64(tt.amount) / float64(tt.bb)
				if tt.bb == 0 {
					actual = float64(tt.amount) / 20
				}
				t.Errorf("%s: amount=%d, bb=%d (ratio=%.2f): expected %v, got %v",
					tt.note, tt.amount, tt.bb, actual, tt.expected, result)
			}
		})
	}
}

// BOUNDARY VALUE TEST SUITE for bucketByPotFraction
// Critical boundaries: 0, 0.38, 0.58, 0.78, 1.15
func TestBucketByPotFraction(t *testing.T) {
	tests := []struct {
		name     string
		amount   int
		pot      int
		expected RangeActionBucket
		note     string
	}{
		// Lower boundary (0)
		{"0 (lower)", 0, 100, RangeActionCheck, "min"},

		// Boundary at 0.38: amount/pot <= 0.38 → BetSmall
		{"exactly 0.38", 38, 100, RangeActionBetSmall, "at-boundary"},
		{"0.37 (37/100)", 37, 100, RangeActionBetSmall, "just-below"},
		{"0.39 (39/100)", 39, 100, RangeActionBetHalf, "just-above"}, // 0.39 > 0.38

		// Boundary at 0.58: amount/pot <= 0.58 → BetHalf
		{"exactly 0.58", 58, 100, RangeActionBetHalf, "at-boundary"},
		{"0.57 (57/100)", 57, 100, RangeActionBetHalf, "just-below"},
		{"0.59 (59/100)", 59, 100, RangeActionBetTwoThird, "just-above"}, // 0.59 > 0.58

		// Boundary at 0.78: amount/pot <= 0.78 → BetTwoThird
		{"exactly 0.78", 78, 100, RangeActionBetTwoThird, "at-boundary"},
		{"0.77 (77/100)", 77, 100, RangeActionBetTwoThird, "just-below"},
		{"0.79 (79/100)", 79, 100, RangeActionBetPot, "just-above"}, // 0.79 > 0.78

		// Boundary at 1.15: amount/pot <= 1.15 → BetPot
		{"exactly 1.15", 115, 100, RangeActionBetPot, "at-boundary"},
		{"1.14 (114/100)", 114, 100, RangeActionBetPot, "just-below"},
		{"1.16 (116/100)", 116, 100, RangeActionBetOver, "just-above"}, // 1.16 > 1.15

		// Upper boundary (very large)
		{"10.0 (extreme)", 1000, 100, RangeActionBetOver, "max"},

		// Edge case: zero pot defaults to BetHalf
		{"zero pot edge case", 50, 0, RangeActionBetHalf, "default-pot"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := &parser.Hand{TotalPot: tt.pot}
			result := bucketByPotFraction(tt.amount, hand)
			if result != tt.expected {
				actual := 0.0
				if tt.pot > 0 {
					actual = float64(tt.amount) / float64(tt.pot)
				}
				t.Errorf("%s: amount=%d, pot=%d (ratio=%.2f): expected %v, got %v",
					tt.note, tt.amount, tt.pot, actual, tt.expected, result)
			}
		})
	}
}

// BOUNDARY VALUE TEST SUITE for preflopRangeActionSummary
func TestPreflopRangeActionSummary(t *testing.T) {
	tests := []struct {
		name     string
		hand     *parser.Hand
		player   *parser.PlayerHandInfo
		expected RangeActionBucket
		ok       bool
		note     string
	}{
		// Boundary: nil player
		{"nil player", &parser.Hand{}, nil, RangeActionCheck, false, "nil"},

		// Boundary: player with no actions
		{"empty actions", &parser.Hand{}, &parser.PlayerHandInfo{}, RangeActionCheck, false, "empty"},

		// Boundary: only preflop actions
		{"fold PF only", &parser.Hand{}, &parser.PlayerHandInfo{
			FoldedPF: true,
			Actions:  []parser.PlayerAction{{Action: parser.ActionFold, Street: parser.StreetPreFlop}},
		}, RangeActionFold, true, "fold"},

		// Boundary: VPIP without explicit action
		{"VPIP no actions", &parser.Hand{}, &parser.PlayerHandInfo{
			VPIP: true,
		}, RangeActionCall, true, "vpip"},

		// Boundary: ThreeBet without explicit action
		{"3Bet no actions", &parser.Hand{}, &parser.PlayerHandInfo{
			ThreeBet: true,
		}, RangeActionBetHalf, true, "3bet"},

		// Boundary: Only postflop actions (should be ignored)
		{"postflop only", &parser.Hand{}, &parser.PlayerHandInfo{
			Actions: []parser.PlayerAction{
				{Action: parser.ActionCall, Amount: 100, Street: parser.StreetFlop},
			},
		}, RangeActionCheck, false, "no-pf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := preflopRangeActionSummary(tt.hand, tt.player)
			if ok != tt.ok {
				t.Errorf("%s: expected ok=%v, got %v", tt.note, tt.ok, ok)
			}
			if result != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.note, tt.expected, result)
			}
		})
	}
}

// BOUNDARY VALUE TEST SUITE for overallActionSummary
func TestOverallActionSummary(t *testing.T) {
	tests := []struct {
		name     string
		hand     *parser.Hand
		player   *parser.PlayerHandInfo
		expected RangeActionBucket
		ok       bool
		note     string
	}{
		// Boundary: nil player
		{"nil player", &parser.Hand{}, nil, RangeActionCheck, false, "nil"},

		// Boundary: no actions
		{"empty actions", &parser.Hand{}, &parser.PlayerHandInfo{}, RangeActionCheck, false, "empty"},

		// Boundary: showdown actions ARE counted (not filtered by street)
		{"showdown only", &parser.Hand{}, &parser.PlayerHandInfo{
			Actions: []parser.PlayerAction{{Action: parser.ActionCall, Street: parser.StreetShowdown}},
		}, RangeActionCall, true, "showdown"},

		// Boundary: single action of each type
		{"single fold", &parser.Hand{}, &parser.PlayerHandInfo{
			Actions: []parser.PlayerAction{{Action: parser.ActionFold, Street: parser.StreetFlop}},
		}, RangeActionFold, true, "fold"},

		{"single check", &parser.Hand{}, &parser.PlayerHandInfo{
			Actions: []parser.PlayerAction{{Action: parser.ActionCheck, Street: parser.StreetFlop}},
		}, RangeActionCheck, true, "check"},

		{"single call", &parser.Hand{}, &parser.PlayerHandInfo{
			Actions: []parser.PlayerAction{{Action: parser.ActionCall, Street: parser.StreetFlop}},
		}, RangeActionCall, true, "call"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := overallActionSummary(tt.hand, tt.player)
			if ok != tt.ok {
				t.Errorf("%s: expected ok=%v, got %v", tt.note, tt.ok, ok)
			}
			if result != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.note, tt.expected, result)
			}
		})
	}
}

// BOUNDARY VALUE TEST SUITE for newHandRangeTable
func TestNewHandRangeTable(t *testing.T) {
	table := newHandRangeTable()

	if table == nil {
		t.Fatal("newHandRangeTable returned nil")
	}

	// Boundary: 13x13 grid must be complete
	if len(table.Cells) != 13 {
		t.Errorf("boundary - row count: expected 13, got %d", len(table.Cells))
	}
	for i := 0; i < 13; i++ {
		if len(table.Cells[i]) != 13 {
			t.Errorf("boundary - col count at row %d: expected 13, got %d", i, len(table.Cells[i]))
		}
	}

	// Boundary: Corner cells
	corners := []struct{ r, c int }{
		{0, 0}, {0, 12}, {12, 0}, {12, 12},
	}
	for _, corner := range corners {
		if table.Cells[corner.r][corner.c] == nil {
			t.Errorf("boundary - corner[%d,%d] is nil", corner.r, corner.c)
		}
	}

	// Boundary: Diagonal cells are pairs
	for i := 0; i < 13; i++ {
		if !table.Cells[i][i].IsPair {
			t.Errorf("boundary - diagonal[%d,%d] should be pair", i, i)
		}
	}

	// Boundary: Upper triangle are suited
	for i := 0; i < 13; i++ {
		for j := i + 1; j < 13; j++ {
			if !table.Cells[i][j].Suited {
				t.Errorf("boundary - upper[%d,%d] should be suited", i, j)
			}
		}
	}

	// Boundary: ByHandClass map exists
	if table.ByHandClass == nil {
		t.Error("boundary - ByHandClass map is nil")
	}
}

// Benchmark tests
func BenchmarkBucketByBBMultiple(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bucketByBBMultiple(100, 20)
	}
}

func BenchmarkBucketByPotFraction(b *testing.B) {
	hand := &parser.Hand{TotalPot: 200}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bucketByPotFraction(100, hand)
	}
}
