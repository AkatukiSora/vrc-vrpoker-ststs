package stats

import "github.com/AkatukiSora/vrc-vrpoker-ststs/internal/parser"

// PocketCategory classifies a two-card preflop holding.
type PocketCategory int

const (
	PocketPremium         PocketCategory = iota // AA, KK, QQ, AKs, AKo
	PocketSecondPremium                         // JJ, TT, AQs, AQo, AJs, KQs, KQo
	PocketPair                                  // 22-AA (any pocket pair)
	PocketSuitedConnector                       // same suit, rank difference == 1
	PocketSuitedOneGapper                       // same suit, rank difference == 2
	PocketSuited                                // same suit (any)
	PocketAx                                    // contains an Ace
	PocketKx                                    // contains a King
	PocketBroadwayOffsuit                       // both cards T-A, different suits
	PocketConnector                             // rank difference == 1 (any suit)
)

// ClassifyPocketHand returns all matching categories for a two-card preflop hand.
// Multiple categories can apply (e.g., AKs is both PocketPremium and PocketSuited).
func ClassifyPocketHand(c1, c2 parser.Card) []PocketCategory {
	r1 := rankValue(c1.Rank)
	r2 := rankValue(c2.Rank)

	diff := r1 - r2
	if diff < 0 {
		diff = -diff
	}

	suited := c1.Suit == c2.Suit
	isPair := r1 == r2
	isBroadway1 := r1 >= 10
	isBroadway2 := r2 >= 10

	var cats []PocketCategory

	// PocketPremium: AA, KK, QQ, or AK (suited or offsuit)
	if isPair && (r1 == 14 || r1 == 13 || r1 == 12) {
		cats = append(cats, PocketPremium)
	} else if (r1 == 14 || r2 == 14) && (r1 == 13 || r2 == 13) {
		cats = append(cats, PocketPremium)
	}

	// PocketSecondPremium: JJ, TT, AQ (suited or offsuit), AJs, KQ (suited or offsuit)
	if isPair && (r1 == 11 || r1 == 10) {
		cats = append(cats, PocketSecondPremium)
	} else if (r1 == 14 || r2 == 14) && (r1 == 12 || r2 == 12) {
		// AQ
		cats = append(cats, PocketSecondPremium)
	} else if suited && (r1 == 14 || r2 == 14) && (r1 == 11 || r2 == 11) {
		// AJs
		cats = append(cats, PocketSecondPremium)
	} else if (r1 == 13 || r2 == 13) && (r1 == 12 || r2 == 12) {
		// KQ
		cats = append(cats, PocketSecondPremium)
	}

	// PocketPair: both cards have the same rank
	if isPair {
		cats = append(cats, PocketPair)
	}

	// PocketSuitedConnector: same suit AND rank diff == 1
	if suited && diff == 1 {
		cats = append(cats, PocketSuitedConnector)
	}

	// PocketSuitedOneGapper: same suit AND rank diff == 2
	if suited && diff == 2 {
		cats = append(cats, PocketSuitedOneGapper)
	}

	// PocketSuited: same suit (regardless of ranks)
	if suited {
		cats = append(cats, PocketSuited)
	}

	// PocketAx: at least one card is Ace
	if r1 == 14 || r2 == 14 {
		cats = append(cats, PocketAx)
	}

	// PocketKx: at least one card is King
	if r1 == 13 || r2 == 13 {
		cats = append(cats, PocketKx)
	}

	// PocketBroadwayOffsuit: both cards Broadway (T+) AND suits differ AND not a pair
	if isBroadway1 && isBroadway2 && !suited && !isPair {
		cats = append(cats, PocketBroadwayOffsuit)
	}

	// PocketConnector: rank diff == 1 (any suit, includes suited connectors)
	if diff == 1 {
		cats = append(cats, PocketConnector)
	}

	return cats
}

// PocketCategoryLabel returns a short display label for a PocketCategory.
func PocketCategoryLabel(c PocketCategory) string {
	switch c {
	case PocketPremium:
		return "Premium"
	case PocketSecondPremium:
		return "2nd Premium"
	case PocketPair:
		return "Pocket Pair"
	case PocketSuitedConnector:
		return "Suited Connector"
	case PocketSuitedOneGapper:
		return "Suited 1-Gap"
	case PocketSuited:
		return "Suited"
	case PocketAx:
		return "Ax"
	case PocketKx:
		return "Kx"
	case PocketBroadwayOffsuit:
		return "Broadway Offsuit"
	case PocketConnector:
		return "Connector"
	default:
		return ""
	}
}

// AllPocketCategories returns all PocketCategory values in display order.
func AllPocketCategories() []PocketCategory {
	return []PocketCategory{
		PocketPremium,
		PocketSecondPremium,
		PocketPair,
		PocketSuitedConnector,
		PocketSuitedOneGapper,
		PocketSuited,
		PocketAx,
		PocketKx,
		PocketBroadwayOffsuit,
		PocketConnector,
	}
}
