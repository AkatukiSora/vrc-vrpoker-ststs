package stats

type MetricID string

const (
	MetricVPIP            MetricID = "vpip"
	MetricPFR             MetricID = "pfr"
	MetricGap             MetricID = "gap"
	MetricRFI             MetricID = "rfi"
	MetricThreeBet        MetricID = "three_bet"
	MetricFourBet         MetricID = "four_bet"
	MetricSqueeze         MetricID = "squeeze"
	MetricFoldToThreeBet  MetricID = "fold_to_three_bet"
	MetricSteal           MetricID = "steal"
	MetricFoldToSteal     MetricID = "fold_to_steal"
	MetricFoldBBToSteal   MetricID = "fold_bb_to_steal"
	MetricFoldSBToSteal   MetricID = "fold_sb_to_steal"
	MetricThreeBetVsSteal MetricID = "three_bet_vs_steal"
	MetricFlopCBet        MetricID = "flop_cbet"
	MetricTurnCBet        MetricID = "turn_cbet"
	MetricFoldToFlopCBet  MetricID = "fold_to_flop_cbet"
	MetricFoldToTurnCBet  MetricID = "fold_to_turn_cbet"
	MetricWTSD            MetricID = "wtsd"
	MetricWSD             MetricID = "w_sd"
	MetricWWSF            MetricID = "wwsf"
	MetricAFq             MetricID = "afq"
	MetricAF              MetricID = "af"
	MetricDelayedCBet     MetricID = "delayed_cbet"
	MetricWonWithoutSD    MetricID = "won_without_showdown"
	MetricBBPer100        MetricID = "bb_per_100"
)

type MetricSampleClass int

const (
	SampleClassHands MetricSampleClass = iota
	SampleClassSituational
)

type MetricFormat int

const (
	MetricFormatPercent MetricFormat = iota
	MetricFormatRatio
	MetricFormatBBPer100
	MetricFormatDiff
)

type MetricDefinition struct {
	ID          MetricID
	Label       string
	SampleClass MetricSampleClass
	Format      MetricFormat
}

type MetricValue struct {
	ID          MetricID
	Count       int
	Opportunity int
	Rate        float64
	Confident   bool
	MinSample   int
	Format      MetricFormat
}

const (
	handFrequencyThreshold = 200
	situationalThreshold   = 50
)

var metricRegistry = []MetricDefinition{
	{ID: MetricVPIP, Label: "VPIP", SampleClass: SampleClassHands, Format: MetricFormatPercent},
	{ID: MetricPFR, Label: "PFR", SampleClass: SampleClassHands, Format: MetricFormatPercent},
	{ID: MetricGap, Label: "Gap", SampleClass: SampleClassHands, Format: MetricFormatDiff},
	{ID: MetricRFI, Label: "RFI", SampleClass: SampleClassHands, Format: MetricFormatPercent},
	{ID: MetricThreeBet, Label: "3Bet", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricFourBet, Label: "4Bet", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricSqueeze, Label: "Squeeze", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricFoldToThreeBet, Label: "Fold to 3Bet", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricSteal, Label: "Steal", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricFoldToSteal, Label: "Fold to Steal", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricFoldBBToSteal, Label: "Fold BB to Steal", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricFoldSBToSteal, Label: "Fold SB to Steal", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricThreeBetVsSteal, Label: "3Bet vs Steal", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricFlopCBet, Label: "Flop CBet", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricTurnCBet, Label: "Turn CBet", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricFoldToFlopCBet, Label: "Fold to Flop CBet", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricFoldToTurnCBet, Label: "Fold to Turn CBet", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricWTSD, Label: "WTSD", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricWSD, Label: "W$SD", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricWWSF, Label: "WWSF", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricAFq, Label: "AFq", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricAF, Label: "AF", SampleClass: SampleClassSituational, Format: MetricFormatRatio},
	{ID: MetricDelayedCBet, Label: "Delayed CBet", SampleClass: SampleClassSituational, Format: MetricFormatPercent},
	{ID: MetricWonWithoutSD, Label: "Won w/o SD", SampleClass: SampleClassHands, Format: MetricFormatPercent},
	{ID: MetricBBPer100, Label: "bb/100", SampleClass: SampleClassHands, Format: MetricFormatBBPer100},
}

func confidenceThreshold(class MetricSampleClass) int {
	if class == SampleClassSituational {
		return situationalThreshold
	}
	return handFrequencyThreshold
}
