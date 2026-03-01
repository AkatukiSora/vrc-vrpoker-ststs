package stats

import (
	"testing"

	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/parser"
)

func TestNewIncrementalCalculator(t *testing.T) {
	ic := NewIncrementalCalculator(0)

	if ic == nil {
		t.Fatal("NewIncrementalCalculator returned nil")
	}
	if ic.localSeat != 0 {
		t.Errorf("expected localSeat=0, got %d", ic.localSeat)
	}
	if ic.s == nil {
		t.Error("stats pointer is nil")
	}
	if ic.ma == nil {
		t.Error("accumulator pointer is nil")
	}
	if ic.calc == nil {
		t.Error("calculator pointer is nil")
	}
	if ic.s.ByPosition == nil {
		t.Error("ByPosition map is nil")
	}
	if ic.s.HandRange == nil {
		t.Error("HandRange is nil")
	}
	if ic.s.Metrics == nil {
		t.Error("Metrics map is nil")
	}
}

func TestIncrementalCalculatorFeedNilHand(t *testing.T) {
	ic := NewIncrementalCalculator(0)

	// Feed nil hand should not panic
	ic.Feed(nil)

	if ic.s.TotalHands != 0 {
		t.Errorf("expected 0 hands after feeding nil, got %d", ic.s.TotalHands)
	}
}

func TestIncrementalCalculatorFeedIncompleteHand(t *testing.T) {
	ic := NewIncrementalCalculator(0)

	hand := &parser.Hand{
		LocalPlayerSeat: 0,
		IsComplete:      false, // Incomplete hand
		Players:         make(map[int]*parser.PlayerHandInfo),
	}
	hand.Players[0] = &parser.PlayerHandInfo{
		Position: parser.PosBTN,
	}

	ic.Feed(hand)

	if ic.s.TotalHands != 0 {
		t.Errorf("expected 0 hands after feeding incomplete hand, got %d", ic.s.TotalHands)
	}
}

func TestIncrementalCalculatorFeedValidHand(t *testing.T) {
	ic := NewIncrementalCalculator(0)

	hand := createValidTestHand(0)
	ic.Feed(hand)

	if ic.s.TotalHands != 1 {
		t.Errorf("expected 1 hand, got %d", ic.s.TotalHands)
	}
}

func TestIncrementalCalculatorFeedMultipleHands(t *testing.T) {
	ic := NewIncrementalCalculator(0)

	hand1 := createValidTestHand(0)
	hand2 := createValidTestHand(0)
	hand3 := createValidTestHand(0)

	ic.Feed(hand1)
	ic.Feed(hand2)
	ic.Feed(hand3)

	if ic.s.TotalHands != 3 {
		t.Errorf("expected 3 hands, got %d", ic.s.TotalHands)
	}
}

func TestIncrementalCalculatorHandCount(t *testing.T) {
	ic := NewIncrementalCalculator(0)

	if ic.HandCount() != 0 {
		t.Errorf("expected 0 hands initially, got %d", ic.HandCount())
	}

	ic.Feed(createValidTestHand(0))
	if ic.HandCount() != 1 {
		t.Errorf("expected 1 hand after feed, got %d", ic.HandCount())
	}

	ic.Feed(createValidTestHand(0))
	if ic.HandCount() != 2 {
		t.Errorf("expected 2 hands after second feed, got %d", ic.HandCount())
	}
}

func TestIncrementalCalculatorCompute(t *testing.T) {
	ic := NewIncrementalCalculator(0)

	hand := createValidTestHand(0)
	ic.Feed(hand)

	stats := ic.Compute()

	if stats == nil {
		t.Fatal("Compute returned nil")
	}
	if stats.TotalHands != 1 {
		t.Errorf("expected 1 hand in computed stats, got %d", stats.TotalHands)
	}
	// Verify that metrics were computed
	if stats.Metrics == nil {
		t.Error("Metrics map is nil in computed stats")
	}
}

func TestIncrementalCalculatorComputeDoesNotMutateAccumulator(t *testing.T) {
	ic := NewIncrementalCalculator(0)

	hand := createValidTestHand(0)
	ic.Feed(hand)

	// Get original state
	originalHandCount := ic.HandCount()

	// Call Compute
	stats1 := ic.Compute()

	// Verify accumulator is unchanged
	if ic.HandCount() != originalHandCount {
		t.Error("Compute mutated the accumulator")
	}

	// Call Compute again and verify it returns same state
	stats2 := ic.Compute()

	if stats2.TotalHands != stats1.TotalHands {
		t.Error("Multiple Compute calls should return consistent results")
	}
}

func TestIncrementalCalculatorComputeAfterFeedUpdates(t *testing.T) {
	ic := NewIncrementalCalculator(0)

	hand1 := createValidTestHand(0)
	ic.Feed(hand1)
	stats1 := ic.Compute()

	if stats1.TotalHands != 1 {
		t.Errorf("expected 1 hand in first compute, got %d", stats1.TotalHands)
	}

	hand2 := createValidTestHand(0)
	ic.Feed(hand2)
	stats2 := ic.Compute()

	if stats2.TotalHands != 2 {
		t.Errorf("expected 2 hands in second compute, got %d", stats2.TotalHands)
	}
}

func TestClonePositionStatsEmpty(t *testing.T) {
	original := make(map[parser.Position]*PositionStats)
	cloned := clonePositionStats(original)

	if cloned == nil {
		t.Error("clonePositionStats returned nil for empty map")
	}
	if len(cloned) != 0 {
		t.Errorf("expected empty cloned map, got %d entries", len(cloned))
	}
}

func TestClonePositionStats(t *testing.T) {
	original := make(map[parser.Position]*PositionStats)
	original[parser.PosBTN] = &PositionStats{
		Position: parser.PosBTN,
		Hands:    10,
		Won:      5,
	}
	original[parser.PosSB] = &PositionStats{
		Position: parser.PosSB,
		Hands:    8,
		Won:      3,
	}

	cloned := clonePositionStats(original)

	// Verify all entries are cloned
	if len(cloned) != 2 {
		t.Errorf("expected 2 entries, got %d", len(cloned))
	}

	// Verify values are correct
	if cloned[parser.PosBTN].Hands != 10 {
		t.Error("cloned BTN stats incorrect")
	}
	if cloned[parser.PosSB].Hands != 8 {
		t.Error("cloned SB stats incorrect")
	}

	// Verify maps are independent
	original[parser.PosBTN].Hands = 20
	if cloned[parser.PosBTN].Hands != 10 {
		t.Error("modifying original affected clone")
	}
}

func TestCloneHandRangeTableNil(t *testing.T) {
	cloned := cloneHandRangeTable(nil)
	if cloned != nil {
		t.Error("cloneHandRangeTable should return nil for nil input")
	}
}

func TestCloneHandRangeTable(t *testing.T) {
	original := newHandRangeTable()

	// Modify some cells to test cloning
	original.Cells[0][0].Dealt = 10
	original.Cells[0][0].Won = 5

	cloned := cloneHandRangeTable(original)

	if cloned == nil {
		t.Fatal("cloneHandRangeTable returned nil")
	}

	// Verify grid structure
	if len(cloned.Cells) != 13 {
		t.Errorf("expected 13 rows, got %d", len(cloned.Cells))
	}

	// Verify cell values are cloned
	if cloned.Cells[0][0].Dealt != 10 {
		t.Error("cell dealt value not cloned")
	}
	if cloned.Cells[0][0].Won != 5 {
		t.Error("cell won value not cloned")
	}

	// Verify maps are independent
	original.Cells[0][0].Dealt = 20
	if cloned.Cells[0][0].Dealt != 10 {
		t.Error("modifying original affected clone")
	}
}

func TestIncrementalCalculatorFeedWithDifferentLocalSeat(t *testing.T) {
	ic := NewIncrementalCalculator(1) // Local seat is 1

	hand := &parser.Hand{
		LocalPlayerSeat: 1,
		IsComplete:      true,
		StatsEligible:   true,
		Players:         make(map[int]*parser.PlayerHandInfo),
	}
	hand.Players[1] = &parser.PlayerHandInfo{
		Position:     parser.PosBTN,
		Participated: true,
		Won:          true,
		PotWon:       100,
	}

	ic.Feed(hand)

	if ic.s.TotalHands != 1 {
		t.Error("should accept hand when local seat matches")
	}
}

// Helper function to create a valid test hand
func createValidTestHand(localSeat int) *parser.Hand {
	hand := &parser.Hand{
		ID:              1,
		LocalPlayerSeat: localSeat,
		NumPlayers:      4,
		SBSeat:          1,
		BBSeat:          2,
		Players:         make(map[int]*parser.PlayerHandInfo),
		CommunityCards: []parser.Card{
			{Rank: "A", Suit: "h"},
			{Rank: "K", Suit: "d"},
			{Rank: "Q", Suit: "c"},
		},
		TotalPot:      100,
		IsComplete:    true,
		StatsEligible: true,
	}

	// Add local player
	hand.Players[localSeat] = &parser.PlayerHandInfo{
		SeatID:       localSeat,
		Position:     parser.PosBTN,
		HoleCards:    []parser.Card{{Rank: "A", Suit: "h"}, {Rank: "K", Suit: "d"}},
		VPIP:         true,
		Participated: true,
		Won:          true,
		PotWon:       100,
		ShowedDown:   false,
		Actions: []parser.PlayerAction{
			{Action: parser.ActionCall, Amount: 20, Street: parser.StreetPreFlop},
		},
	}

	// Mark hand as stats-eligible (it will check if IsComplete and other conditions)
	return hand
}

// Benchmark tests
func BenchmarkNewIncrementalCalculator(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewIncrementalCalculator(0)
	}
}

func BenchmarkIncrementalCalculatorFeed(b *testing.B) {
	ic := NewIncrementalCalculator(0)
	hand := createValidTestHand(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ic.Feed(hand)
	}
}

func BenchmarkIncrementalCalculatorCompute(b *testing.B) {
	ic := NewIncrementalCalculator(0)
	for i := 0; i < 100; i++ {
		ic.Feed(createValidTestHand(0))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ic.Compute()
	}
}

func BenchmarkClonePositionStats(b *testing.B) {
	original := make(map[parser.Position]*PositionStats)
	original[parser.PosBTN] = &PositionStats{Position: parser.PosBTN, Hands: 100}
	original[parser.PosSB] = &PositionStats{Position: parser.PosSB, Hands: 50}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clonePositionStats(original)
	}
}

func BenchmarkCloneHandRangeTable(b *testing.B) {
	original := newHandRangeTable()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cloneHandRangeTable(original)
	}
}
