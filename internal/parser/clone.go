package parser

func (p *Parser) Clone() *Parser {
	if p == nil {
		return nil
	}

	clone := &Parser{
		result:           cloneParseResult(p.result),
		handIDCounter:    p.handIDCounter,
		inPokerWorld:     p.inPokerWorld,
		worldDetected:    p.worldDetected,
		currentStreet:    p.currentStreet,
		streetBetAmount:  p.streetBetAmount,
		streetBets:       cloneIntMap(p.streetBets),
		foldedThisHand:   cloneBoolMap(p.foldedThisHand),
		pendingWinners:   append([]pendingWin(nil), p.pendingWinners...),
		lastTimestamp:    p.lastTimestamp,
		pendingLocalSeat: p.pendingLocalSeat,
		lastBlindSeat:    p.lastBlindSeat,
		pfActions:        append([]pfAction(nil), p.pfActions...),
	}
	clone.pendingLocalCards = append([]Card(nil), p.pendingLocalCards...)
	clone.currentHand = cloneHand(p.currentHand)
	if clone.currentHand != nil {
		clone.result.CurrentHand = clone.currentHand
	}
	return clone
}

func cloneParseResult(in ParseResult) ParseResult {
	out := in
	out.Hands = make([]*Hand, 0, len(in.Hands))
	for _, hand := range in.Hands {
		out.Hands = append(out.Hands, cloneHand(hand))
	}
	if in.CurrentHand != nil {
		out.CurrentHand = cloneHand(in.CurrentHand)
	}
	return out
}

func cloneHand(h *Hand) *Hand {
	if h == nil {
		return nil
	}
	copyHand := *h
	copyHand.CommunityCards = append([]Card(nil), h.CommunityCards...)
	copyHand.ActiveSeats = append([]int(nil), h.ActiveSeats...)
	copyHand.Players = make(map[int]*PlayerHandInfo, len(h.Players))
	for seat, pi := range h.Players {
		copyHand.Players[seat] = clonePlayerInfo(pi)
	}
	return &copyHand
}

func clonePlayerInfo(pi *PlayerHandInfo) *PlayerHandInfo {
	if pi == nil {
		return nil
	}
	copyPI := *pi
	copyPI.HoleCards = append([]Card(nil), pi.HoleCards...)
	copyPI.Actions = append([]PlayerAction(nil), pi.Actions...)
	return &copyPI
}

func cloneIntMap(in map[int]int) map[int]int {
	if len(in) == 0 {
		return make(map[int]int)
	}
	out := make(map[int]int, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneBoolMap(in map[int]bool) map[int]bool {
	if len(in) == 0 {
		return make(map[int]bool)
	}
	out := make(map[int]bool, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
