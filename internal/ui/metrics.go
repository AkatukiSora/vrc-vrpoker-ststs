package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"

	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/stats"
)

type MetricValue struct {
	Display       string
	Color         color.Color
	Opportunities int
}

type MetricDefinition struct {
	ID             string
	Label          string
	HelpKey        string
	HelpFallback   string
	MinSamples     int
	GoodSamples    int
	ShowInOverview bool
	ShowInPosition bool
	StatsMetricID  *stats.MetricID
	OverviewValue  func(*stats.Stats) MetricValue
	PositionValue  func(*stats.PositionStats) MetricValue
}

func (m MetricDefinition) HelpText() string {
	if m.HelpKey == "" {
		return m.HelpFallback
	}
	return lang.X(m.HelpKey, m.HelpFallback)
}

type MetricPreset struct {
	Name       string
	MetricIDs  map[string]struct{}
	ButtonText string
}

type MetricVisibilityState struct {
	Visible map[string]bool
}

func NewMetricVisibilityState() *MetricVisibilityState {
	visible := make(map[string]bool, len(metricRegistry))
	for _, m := range metricRegistry {
		visible[m.ID] = true
	}
	return &MetricVisibilityState{Visible: visible}
}

func (m *MetricVisibilityState) IsVisible(metricID string) bool {
	if m == nil {
		return true
	}
	v, ok := m.Visible[metricID]
	if !ok {
		return true
	}
	return v
}

func (m *MetricVisibilityState) SetVisible(metricID string, visible bool) {
	if m == nil {
		return
	}
	if m.Visible == nil {
		m.Visible = make(map[string]bool)
	}
	m.Visible[metricID] = visible
}

func (m *MetricVisibilityState) ApplyPreset(p MetricPreset) {
	if m == nil {
		return
	}
	if m.Visible == nil {
		m.Visible = make(map[string]bool)
	}
	for _, metric := range metricRegistry {
		_, keep := p.MetricIDs[metric.ID]
		m.Visible[metric.ID] = keep
	}
}

func metricsForOverview(visibility *MetricVisibilityState) []MetricDefinition {
	out := make([]MetricDefinition, 0, len(metricRegistry))
	for _, m := range metricRegistry {
		if !m.ShowInOverview {
			continue
		}
		if visibility != nil && !visibility.IsVisible(m.ID) {
			continue
		}
		out = append(out, m)
	}
	return out
}

func metricsForPosition(visibility *MetricVisibilityState) []MetricDefinition {
	out := make([]MetricDefinition, 0, len(metricRegistry))
	for _, m := range metricRegistry {
		if !m.ShowInPosition {
			continue
		}
		if visibility != nil && !visibility.IsVisible(m.ID) {
			continue
		}
		out = append(out, m)
	}
	return out
}

func metricFootnoteText(opportunities, minSamples int) string {
	if opportunities < 0 {
		return ""
	}
	return lang.X("metric.footnote.normal", "n={{.N}}", map[string]any{"N": opportunities})
}

type sampleThreshold struct {
	Min  int
	Good int
}

func thresholdForMetricID(id string) sampleThreshold {
	th := metricThresholdForMetricID(id)
	return sampleThreshold{Min: th.Min, Good: th.Good}
}

func metricPresets() []MetricPreset {
	beginner := setOf(
		"hands", "profit",
		"vpip", "pfr", "gap", "rfi", "steal",
		"three_bet", "fold_to_three_bet", "fold_bb_to_steal", "fold_sb_to_steal",
		"flop_cbet", "turn_cbet", "fold_to_flop_cbet", "fold_to_turn_cbet",
		"wtsd", "w_sd", "wwsf",
	)
	advanced := setOf(
		"hands", "profit", "bb_per_100",
		"vpip", "pfr", "gap", "rfi", "steal",
		"three_bet", "three_bet_vs_steal", "fold_to_three_bet", "four_bet", "squeeze",
		"fold_to_steal", "fold_bb_to_steal", "fold_sb_to_steal",
		"flop_cbet", "turn_cbet", "delayed_cbet", "fold_to_flop_cbet", "fold_to_turn_cbet",
		"wtsd", "w_sd", "wwsf", "afq", "af", "won_without_showdown",
	)
	leak := setOf(
		"vpip", "pfr", "gap", "rfi", "steal",
		"three_bet", "three_bet_vs_steal", "fold_to_three_bet", "four_bet", "squeeze",
		"fold_bb_to_steal", "fold_sb_to_steal", "fold_to_steal",
		"flop_cbet", "turn_cbet", "fold_to_flop_cbet", "fold_to_turn_cbet",
		"wtsd", "w_sd", "wwsf", "afq", "won_without_showdown", "bb_per_100",
	)
	all := setOf()
	for _, m := range metricRegistry {
		all[m.ID] = struct{}{}
	}
	return []MetricPreset{
		{Name: "Beginner", ButtonText: lang.X("preset.beginner", "Beginner"), MetricIDs: beginner},
		{Name: "Advanced", ButtonText: lang.X("preset.advanced", "Advanced"), MetricIDs: advanced},
		{Name: "Leak Focus", ButtonText: lang.X("preset.leak_focus", "Leak Focus"), MetricIDs: leak},
		{Name: "All", ButtonText: lang.X("preset.all", "All"), MetricIDs: all},
	}
}

func setOf(ids ...string) map[string]struct{} {
	out := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		out[id] = struct{}{}
	}
	return out
}

func colorForVPIP(rate float64) color.Color {
	if rate > 35 {
		return color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0xff}
	}
	if rate >= 20 {
		return color.NRGBA{R: 0xff, G: 0xd6, B: 0x00, A: 0xff}
	}
	return theme.ForegroundColor()
}

func colorForProfit(profit int) color.Color {
	if profit > 0 {
		return uiSuccessAccent
	}
	if profit < 0 {
		return uiDangerAccent
	}
	return theme.ForegroundColor()
}

func statsMetricValue(s *stats.Stats, id stats.MetricID) MetricValue {
	if s == nil {
		return MetricValue{Display: "-", Color: theme.ForegroundColor(), Opportunities: 0}
	}
	m, ok := s.Metric(id)
	if !ok {
		return MetricValue{Display: "-", Color: theme.ForegroundColor(), Opportunities: 0}
	}
	switch m.Format {
	case stats.MetricFormatRatio:
		return MetricValue{Display: fmt.Sprintf("%.2f", m.Rate), Color: theme.ForegroundColor(), Opportunities: m.Opportunity}
	case stats.MetricFormatBBPer100:
		return MetricValue{Display: fmt.Sprintf("%.2f", m.Rate), Color: theme.ForegroundColor(), Opportunities: m.Opportunity}
	case stats.MetricFormatDiff, stats.MetricFormatPercent:
		return MetricValue{Display: fmt.Sprintf("%.1f%%", m.Rate), Color: theme.ForegroundColor(), Opportunities: m.Opportunity}
	default:
		return MetricValue{Display: fmt.Sprintf("%.1f%%", m.Rate), Color: theme.ForegroundColor(), Opportunities: m.Opportunity}
	}
}

type evidenceItem struct {
	Text   string
	IsGood bool
}

type trendInsight struct {
	Priority  string
	Severity  string
	Title     string
	Text      string
	Evidence  []evidenceItem
	LowSample bool
}

func insightSeverityLabel(code string) string {
	switch code {
	case "P0":
		return lang.X("insight.severity.high", "High")
	case "P1":
		return lang.X("insight.severity.medium", "Medium")
	default:
		return lang.X("insight.severity.low", "Low")
	}
}

func insightEvidenceLine(s *stats.Stats, id stats.MetricID, label, normal string, lo, hi float64, reason, goodReason string) evidenceItem {
	m, ok := s.Metric(id)
	if !ok {
		return evidenceItem{}
	}
	var v string
	if m.Format == stats.MetricFormatRatio || m.Format == stats.MetricFormatBBPer100 {
		v = fmt.Sprintf("%.2f", m.Rate)
	} else {
		v = fmt.Sprintf("%.1f%%", m.Rate)
	}
	n := m.Opportunity
	if m.Rate >= lo && m.Rate <= hi {
		return evidenceItem{
			Text: lang.X("insight.evidence.good_line", "{{.Label}} {{.Value}} (n={{.N}}) | Typical: {{.Normal}} | Good: {{.GoodReason}}", map[string]any{
				"Label":      label,
				"Value":      v,
				"N":          n,
				"Normal":     normal,
				"GoodReason": goodReason,
			}),
			IsGood: true,
		}
	}
	return evidenceItem{
		Text: lang.X("insight.evidence.line", "{{.Label}} {{.Value}} (n={{.N}}) | Typical: {{.Normal}} | {{.Reason}}", map[string]any{
			"Label":  label,
			"Value":  v,
			"N":      n,
			"Normal": normal,
			"Reason": reason,
		}),
		IsGood: false,
	}
}

func buildTrendInsights(s *stats.Stats) []trendInsight {
	if s == nil || s.Metrics == nil {
		return nil
	}
	out := make([]trendInsight, 0, 8)
	add := func(priority, titleKey, titleFallback, textKey, textFallback string, lowSample bool, evidence []evidenceItem) {
		cleanEvidence := make([]evidenceItem, 0, len(evidence))
		for _, e := range evidence {
			if e.Text == "" {
				continue
			}
			cleanEvidence = append(cleanEvidence, e)
		}
		out = append(out, trendInsight{
			Priority:  priority,
			Severity:  insightSeverityLabel(priority),
			Title:     lang.X(titleKey, titleFallback),
			Text:      lang.X(textKey, textFallback),
			Evidence:  cleanEvidence,
			LowSample: lowSample,
		})
	}
	lowBy := func(metricID stats.MetricID) bool {
		m, ok := s.Metric(metricID)
		if !ok {
			return true
		}
		return m.Opportunity < thresholdForMetricID(string(metricID)).Min
	}
	vpip, okV := s.Metric(stats.MetricVPIP)
	pfr, okP := s.Metric(stats.MetricPFR)
	if okV && okP && vpip.Rate-pfr.Rate >= 11 {
		add("P0", "insight.passive_entry.title", "Passive preflop entries", "insight.passive_entry.text", "You may be entering too many pots by call. Consider shifting to a raise-first plan in open spots.", lowBy(stats.MetricVPIP) || lowBy(stats.MetricPFR), []evidenceItem{
			insightEvidenceLine(s, stats.MetricVPIP, "VPIP", "18-28%", 18.0, 28.0, lang.X("insight.reason.vpip_high", "Participation is wider than standard ranges."), lang.X("insight.reason.vpip_good", "Participation rate is within a healthy range.")),
			insightEvidenceLine(s, stats.MetricPFR, "PFR", "12-22%", 12.0, 22.0, lang.X("insight.reason.pfr_low", "Raise frequency is not keeping up with VPIP."), lang.X("insight.reason.pfr_good", "Raise frequency is balanced relative to VPIP.")),
			insightEvidenceLine(s, stats.MetricGap, "Gap", "0-10", 0.0, 10.0, lang.X("insight.reason.gap_high", "Large VPIP-PFR gap suggests passive calls."), lang.X("insight.reason.gap_good", "VPIP-PFR gap is within a healthy range.")),
		})
	}

	threeBet, ok3b := s.Metric(stats.MetricThreeBet)
	foldTo3Bet, okF3 := s.Metric(stats.MetricFoldToThreeBet)
	if ok3b && threeBet.Rate <= 3.5 {
		add("P0", "insight.preflop_exploit.title", "Preflop exploit risk", "insight.preflop_exploit.text", "Low 3-bet frequency can let opponents open too wide against you.", lowBy(stats.MetricThreeBet), []evidenceItem{
			insightEvidenceLine(s, stats.MetricThreeBet, "3Bet", "4-9%", 4.0, 9.0, lang.X("insight.reason.threebet_low", "Too few re-raises allow wider opens."), lang.X("insight.reason.threebet_good", "3-bet frequency is within standard ranges.")),
			insightEvidenceLine(s, stats.MetricThreeBetVsSteal, "3Bet vs Steal", "5-12%", 5.0, 12.0, lang.X("insight.reason.threebet_vs_steal_low", "Blind counter-pressure versus steals is limited."), lang.X("insight.reason.threebet_vs_steal_good", "Counter-pressure versus steals is adequate.")),
		})
	}
	if okF3 && foldTo3Bet.Rate >= 70 {
		add("P0", "insight.fold_to_3bet.title", "Open is too vulnerable to 3-bets", "insight.fold_to_3bet.text", "You fold too often versus 3-bets after opening. Opponents may 3-bet you aggressively.", lowBy(stats.MetricFoldToThreeBet), []evidenceItem{
			insightEvidenceLine(s, stats.MetricFoldToThreeBet, "Fold to 3Bet", "40-55%", 40.0, 55.0, lang.X("insight.reason.fold_to_3bet_high", "This fold rate is high enough to invite aggressive 3-bets."), lang.X("insight.reason.fold_to_3bet_good", "3-bet fold rate is within a balanced range.")),
			insightEvidenceLine(s, stats.MetricFourBet, "4Bet", "1-3%", 1.0, 3.0, lang.X("insight.reason.fourbet_low", "Low 4-bet frequency gives fewer counter options."), lang.X("insight.reason.fourbet_good", "4-bet frequency is within a typical range.")),
		})
	}

	foldBBSteal, okFBB := s.Metric(stats.MetricFoldBBToSteal)
	foldSBSteal, okFSB := s.Metric(stats.MetricFoldSBToSteal)
	if (okFBB && foldBBSteal.Rate >= 65) || (okFSB && foldSBSteal.Rate >= 70) {
		add("P0", "insight.overfold_blinds.title", "Overfolding in blinds", "insight.overfold_blinds.text", "You may be folding too much versus steals, which is easy to exploit over many hands.", (okFBB && lowBy(stats.MetricFoldBBToSteal)) || (okFSB && lowBy(stats.MetricFoldSBToSteal)), []evidenceItem{
			insightEvidenceLine(s, stats.MetricFoldBBToSteal, "Fold BB to Steal", "40-55%", 40.0, 55.0, lang.X("insight.reason.fold_bb_high", "Big blind defense is below a typical defend mix."), lang.X("insight.reason.fold_bb_good", "Big blind fold rate is within a healthy range.")),
			insightEvidenceLine(s, stats.MetricFoldSBToSteal, "Fold SB to Steal", "45-60%", 45.0, 60.0, lang.X("insight.reason.fold_sb_high", "Small blind folds are high versus steals."), lang.X("insight.reason.fold_sb_good", "Small blind fold rate is within a healthy range.")),
		})
	}
	if (okFBB && foldBBSteal.Rate <= 35) || (okFSB && foldSBSteal.Rate <= 35) {
		add("P1", "insight.overdefend_blinds.title", "Over-defending blinds", "insight.overdefend_blinds.text", "You may be defending too wide out of position, leading to difficult postflop spots.", (okFBB && lowBy(stats.MetricFoldBBToSteal)) || (okFSB && lowBy(stats.MetricFoldSBToSteal)), []evidenceItem{
			insightEvidenceLine(s, stats.MetricFoldBBToSteal, "Fold BB to Steal", "40-55%", 40.0, 55.0, lang.X("insight.reason.fold_bb_low", "Very low fold rate can over-expand OOP defense."), lang.X("insight.reason.fold_bb_good", "Big blind fold rate is within a healthy range.")),
			insightEvidenceLine(s, stats.MetricFoldSBToSteal, "Fold SB to Steal", "45-60%", 45.0, 60.0, lang.X("insight.reason.fold_sb_low", "Very low fold rate can over-expand OOP defense."), lang.X("insight.reason.fold_sb_good", "Small blind fold rate is within a healthy range.")),
			insightEvidenceLine(s, stats.MetricWTSD, "WTSD", "22-30%", 22.0, 30.0, lang.X("insight.reason.wtsd_support", "Showdown tendency helps confirm over-calling risk."), lang.X("insight.reason.wtsd_good", "Showdown reach frequency is well-balanced.")),
		})
	}

	rfi, okRFI := s.Metric(stats.MetricRFI)
	steal, okSteal := s.Metric(stats.MetricSteal)
	if (okRFI && rfi.Rate <= 16) || (okSteal && steal.Rate <= 28) {
		add("P1", "insight.missed_steal.title", "Missed steal/value opportunities", "insight.missed_steal.text", "Late-position opens may be too tight. You could be leaving uncontested pots on the table.", (okRFI && lowBy(stats.MetricRFI)) || (okSteal && lowBy(stats.MetricSteal)), []evidenceItem{
			insightEvidenceLine(s, stats.MetricRFI, "RFI", "15-25% (MP), 25-55% (CO/BTN)", 16.0, 55.0, lang.X("insight.reason.rfi_low", "Open frequency is conservative for steal-heavy positions."), lang.X("insight.reason.rfi_good", "Open frequency is within a healthy range.")),
			insightEvidenceLine(s, stats.MetricSteal, "Steal Attempt", "30-50%", 30.0, 50.0, lang.X("insight.reason.steal_low", "Steal spots are not converted often enough."), lang.X("insight.reason.steal_good", "Steal conversion rate is within a healthy range.")),
		})
	}

	foldFlopCBet, okFFC := s.Metric(stats.MetricFoldToFlopCBet)
	foldTurnCBet, okFTC := s.Metric(stats.MetricFoldToTurnCBet)
	if okFFC && foldFlopCBet.Rate >= 60 {
		add("P0", "insight.overfold_flop.title", "Overfolding vs flop c-bets", "insight.overfold_flop.text", "Opponents may profit by c-betting very wide because you fold too frequently on the flop.", lowBy(stats.MetricFoldToFlopCBet), []evidenceItem{
			insightEvidenceLine(s, stats.MetricFoldToFlopCBet, "Fold to Flop CBet", "35-50%", 35.0, 50.0, lang.X("insight.reason.fold_flop_high", "Flop folds are above a defend-balanced range."), lang.X("insight.reason.fold_flop_good", "Flop fold rate is within a balanced defend range.")),
			insightEvidenceLine(s, stats.MetricWWSF, "WWSF", "42-48%", 42.0, 48.0, lang.X("insight.reason.wwsf_low", "Low flop-win frequency supports an overfold pattern."), lang.X("insight.reason.wwsf_good", "Postflop pot capture rate is within a typical range.")),
		})
	}
	if okFTC && foldTurnCBet.Rate >= 65 {
		add("P1", "insight.overfold_turn.title", "Overfolding vs turn barrels", "insight.overfold_turn.text", "You may be giving up too often on turn pressure after defending flop.", lowBy(stats.MetricFoldToTurnCBet), []evidenceItem{
			insightEvidenceLine(s, stats.MetricFoldToFlopCBet, "Fold to Flop CBet", "35-50%", 35.0, 50.0, lang.X("insight.reason.fold_flop_normal_turn_high", "Flop defense is acceptable but turn folds spike."), lang.X("insight.reason.fold_flop_good", "Flop fold rate is within a balanced defend range.")),
			insightEvidenceLine(s, stats.MetricFoldToTurnCBet, "Fold to Turn CBet", "40-55%", 40.0, 55.0, lang.X("insight.reason.fold_turn_high", "Turn folds are high versus typical pressure handling."), lang.X("insight.reason.fold_turn_good", "Turn fold rate is within a healthy range.")),
		})
	}

	flopCbet, okFC := s.Metric(stats.MetricFlopCBet)
	turnCbet, okTC := s.Metric(stats.MetricTurnCBet)
	wwsf, okWW := s.Metric(stats.MetricWWSF)
	if okFC && okTC && okWW && flopCbet.Rate >= 75 && turnCbet.Rate <= 30 {
		add("P1", "insight.auto_cbet.title", "Auto c-bet tendency", "insight.auto_cbet.text", "High flop c-bet with low turn follow-through may indicate one-and-done aggression.", lowBy(stats.MetricFlopCBet) || lowBy(stats.MetricTurnCBet) || lowBy(stats.MetricWWSF), []evidenceItem{
			insightEvidenceLine(s, stats.MetricFlopCBet, "Flop CBet", "50-70%", 50.0, 70.0, lang.X("insight.reason.flop_cbet_high", "Flop c-bet rate is above standard continuation ranges."), lang.X("insight.reason.flop_cbet_good", "Flop c-bet frequency is within standard ranges.")),
			insightEvidenceLine(s, stats.MetricTurnCBet, "Turn CBet", "30-55%", 30.0, 55.0, lang.X("insight.reason.turn_cbet_low", "Turn follow-through is low after flop aggression."), lang.X("insight.reason.turn_cbet_good", "Turn follow-through frequency is adequate.")),
			insightEvidenceLine(s, stats.MetricWWSF, "WWSF", "42-48%", 42.0, 48.0, lang.X("insight.reason.wwsf_support", "Low capture rate supports one-and-done concern."), lang.X("insight.reason.wwsf_good", "Postflop pot capture rate is within a typical range.")),
		})
	}

	wtsd, okWT := s.Metric(stats.MetricWTSD)
	wsd, okWSD := s.Metric(stats.MetricWSD)
	if okWT && okWSD && wtsd.Rate >= 32 && wsd.Rate <= 45 {
		add("P0", "insight.overcall_sd.title", "Over-calling to showdown", "insight.overcall_sd.text", "High WTSD with low W$SD often means too many thin calls in marginal bluff-catch spots.", lowBy(stats.MetricWTSD) || lowBy(stats.MetricWSD), []evidenceItem{
			insightEvidenceLine(s, stats.MetricWTSD, "WTSD", "22-30%", 22.0, 30.0, lang.X("insight.reason.wtsd_high", "Showdown frequency is high for a balanced line."), lang.X("insight.reason.wtsd_good", "Showdown reach frequency is well-balanced.")),
			insightEvidenceLine(s, stats.MetricWSD, "W$SD", "47-55%", 47.0, 55.0, lang.X("insight.reason.wsd_low", "Lower showdown win rate suggests thin calls."), lang.X("insight.reason.wsd_good", "Showdown win rate is within a healthy range.")),
		})
	}
	if okWT && okWW && wtsd.Rate <= 20 && wwsf.Rate < 42 {
		add("P1", "insight.underreach_sd.title", "Not reaching showdown enough", "insight.underreach_sd.text", "Low WTSD with low WWSF can indicate over-folding and missed bluff-catch opportunities.", lowBy(stats.MetricWTSD) || lowBy(stats.MetricWWSF), []evidenceItem{
			insightEvidenceLine(s, stats.MetricWTSD, "WTSD", "22-30%", 22.0, 30.0, lang.X("insight.reason.wtsd_low", "Showdown frequency is low for balanced bluff-catching."), lang.X("insight.reason.wtsd_good", "Showdown reach frequency is well-balanced.")),
			insightEvidenceLine(s, stats.MetricWWSF, "WWSF", "42-48%", 42.0, 48.0, lang.X("insight.reason.wwsf_low", "Postflop pot capture is below typical range."), lang.X("insight.reason.wwsf_good", "Postflop pot capture rate is within a typical range.")),
		})
	}

	afq, okAFq := s.Metric(stats.MetricAFq)
	wonWithoutSD, okWNSD := s.Metric(stats.MetricWonWithoutSD)
	if okWW && wwsf.Rate < 40 {
		add("P0", "insight.low_wwsf.title", "Low postflop pot capture", "insight.low_wwsf.text", "You may be playing too passively postflop and failing to win enough pots after seeing the flop.", lowBy(stats.MetricWWSF), []evidenceItem{
			insightEvidenceLine(s, stats.MetricWWSF, "WWSF", "42-48%", 42.0, 48.0, lang.X("insight.reason.wwsf_low", "Postflop pot capture is below typical range."), lang.X("insight.reason.wwsf_good", "Postflop pot capture rate is within a typical range.")),
			insightEvidenceLine(s, stats.MetricAFq, "AFq", "40-55%", 40.0, 55.0, lang.X("insight.reason.afq_low", "Aggression frequency is on the passive side."), lang.X("insight.reason.afq_good", "Aggression frequency is within a balanced range.")),
			insightEvidenceLine(s, stats.MetricWonWithoutSD, "Won w/o SD", "45-55%", 45.0, 55.0, lang.X("insight.reason.won_without_sd_low", "Non-showdown pot capture is limited."), lang.X("insight.reason.won_without_sd_good", "Non-showdown pot capture rate is within a healthy range.")),
		})
	}
	if okWNSD && wonWithoutSD.Rate > 58 && okWSD && wsd.Rate < 47 && okAFq && afq.Rate >= 55 {
		add("P2", "insight.overbluff_bias.title", "Possible over-bluff bias", "insight.overbluff_bias.text", "Very high non-showdown wins with weaker showdown outcomes may become fragile versus stronger opponents.", lowBy(stats.MetricWonWithoutSD) || lowBy(stats.MetricWSD) || lowBy(stats.MetricAFq), []evidenceItem{
			insightEvidenceLine(s, stats.MetricWonWithoutSD, "Won w/o SD", "45-55%", 45.0, 55.0, lang.X("insight.reason.won_without_sd_high", "Non-showdown wins are unusually high."), lang.X("insight.reason.won_without_sd_good", "Non-showdown pot capture rate is within a healthy range.")),
			insightEvidenceLine(s, stats.MetricWSD, "W$SD", "47-55%", 47.0, 55.0, lang.X("insight.reason.wsd_low", "Showdown performance is below standard range."), lang.X("insight.reason.wsd_good", "Showdown win rate is within a healthy range.")),
			insightEvidenceLine(s, stats.MetricAFq, "AFq", "40-55%", 40.0, 55.0, lang.X("insight.reason.afq_high", "Aggression frequency is very high."), lang.X("insight.reason.afq_good", "Aggression frequency is within a balanced range.")),
		})
	}
	return out
}

var metricRegistry = buildMetricRegistry()

func statsMetricDef(id stats.MetricID, label, helpKey, helpFallback string, showPosition bool) MetricDefinition {
	idCopy := id
	th := thresholdForMetricID(string(id))
	return MetricDefinition{
		ID:             string(id),
		Label:          label,
		HelpKey:        helpKey,
		HelpFallback:   helpFallback,
		MinSamples:     th.Min,
		GoodSamples:    th.Good,
		ShowInOverview: true,
		ShowInPosition: showPosition,
		StatsMetricID:  &idCopy,
		OverviewValue: func(s *stats.Stats) MetricValue {
			v := statsMetricValue(s, id)
			if id == stats.MetricVPIP {
				v.Color = colorForVPIP(vRateOrZero(s, id))
			}
			return v
		},
		PositionValue: func(ps *stats.PositionStats) MetricValue {
			if ps == nil {
				return MetricValue{Display: "-", Color: theme.ForegroundColor(), Opportunities: 0}
			}
			switch id {
			case stats.MetricVPIP:
				rate := ps.VPIPRate()
				return MetricValue{Display: fmt.Sprintf("%.1f%%", rate), Color: colorForVPIP(rate), Opportunities: ps.Hands}
			case stats.MetricPFR:
				return MetricValue{Display: fmt.Sprintf("%.1f%%", ps.PFRRate()), Color: theme.ForegroundColor(), Opportunities: ps.Hands}
			case stats.MetricThreeBet:
				return MetricValue{Display: fmt.Sprintf("%.1f%%", ps.ThreeBetRate()), Color: theme.ForegroundColor(), Opportunities: ps.ThreeBetOpp}
			case stats.MetricFoldToThreeBet:
				return MetricValue{Display: fmt.Sprintf("%.1f%%", ps.FoldTo3BetRate()), Color: theme.ForegroundColor(), Opportunities: ps.FoldTo3BetOpp}
			case stats.MetricWSD:
				return MetricValue{Display: fmt.Sprintf("%.1f%%", ps.WSDRate()), Color: theme.ForegroundColor(), Opportunities: ps.Showdowns}
			default:
				return MetricValue{Display: "-", Color: theme.DisabledColor(), Opportunities: 0}
			}
		},
	}
}

func buildMetricRegistry() []MetricDefinition {
	return []MetricDefinition{
		{
			ID:             "hands",
			Label:          "Hands",
			HelpKey:        "metric.hands.help",
			HelpFallback:   "Total complete hands included in current stats scope.",
			MinSamples:     300,
			GoodSamples:    5000,
			ShowInOverview: true,
			ShowInPosition: true,
			OverviewValue: func(s *stats.Stats) MetricValue {
				if s == nil {
					return MetricValue{Display: "0", Color: theme.ForegroundColor(), Opportunities: 0}
				}
				return MetricValue{Display: fmt.Sprintf("%d", s.TotalHands), Color: theme.ForegroundColor(), Opportunities: s.TotalHands}
			},
			PositionValue: func(ps *stats.PositionStats) MetricValue {
				if ps == nil {
					return MetricValue{Display: "0", Color: theme.ForegroundColor(), Opportunities: 0}
				}
				return MetricValue{Display: fmt.Sprintf("%d", ps.Hands), Color: theme.ForegroundColor(), Opportunities: ps.Hands}
			},
		},
		{
			ID:             "profit",
			Label:          "Total Profit",
			HelpKey:        "metric.profit.help",
			HelpFallback:   "Total chips won minus invested chips.",
			MinSamples:     300,
			GoodSamples:    5000,
			ShowInOverview: true,
			ShowInPosition: true,
			OverviewValue: func(s *stats.Stats) MetricValue {
				if s == nil {
					return MetricValue{Display: "+0", Color: theme.ForegroundColor(), Opportunities: 0}
				}
				p := s.TotalPotWon - s.TotalInvested
				return MetricValue{Display: fmt.Sprintf("%+d", p), Color: colorForProfit(p), Opportunities: s.TotalHands}
			},
			PositionValue: func(ps *stats.PositionStats) MetricValue {
				if ps == nil {
					return MetricValue{Display: "+0", Color: theme.ForegroundColor(), Opportunities: 0}
				}
				p := ps.PotWon - ps.Invested
				return MetricValue{Display: fmt.Sprintf("%+d", p), Color: colorForProfit(p), Opportunities: ps.Hands}
			},
		},
		// Preflop participation and opening
		statsMetricDef(stats.MetricVPIP, "VPIP", "metric.vpip.help", "Voluntarily Put Money In Pot. Preflop participation frequency.", true),
		statsMetricDef(stats.MetricPFR, "PFR", "metric.pfr.help", "Preflop raise frequency.", true),
		statsMetricDef(stats.MetricGap, "VPIP-PFR Gap", "metric.gap.help", "VPIP minus PFR. Larger gap implies more passive preflop entries.", false),
		statsMetricDef(stats.MetricRFI, "RFI", "metric.rfi.help", "Raise First In frequency.", false),
		statsMetricDef(stats.MetricSteal, "Steal Attempt", "metric.steal.help", "Open-raise attempt from steal positions when folded to you.", false),
		// 3-bet/4-bet and preflop pressure
		statsMetricDef(stats.MetricThreeBet, "3Bet", "metric.three_bet.help", "3-bet frequency when a 3-bet opportunity is present.", true),
		statsMetricDef(stats.MetricThreeBetVsSteal, "3Bet vs Steal", "metric.three_bet_vs_steal.help", "3-bet frequency from blinds versus steal opens.", false),
		statsMetricDef(stats.MetricFoldToThreeBet, "Fold to 3Bet", "metric.fold_to_three_bet.help", "Fold frequency when facing a 3-bet after opening.", true),
		statsMetricDef(stats.MetricFourBet, "4Bet", "metric.four_bet.help", "4-bet frequency when facing a 3-bet.", false),
		statsMetricDef(stats.MetricSqueeze, "Squeeze", "metric.squeeze.help", "Squeeze frequency after open + caller before you act.", false),
		// Blind defense
		statsMetricDef(stats.MetricFoldToSteal, "Fold to Steal", "metric.fold_to_steal.help", "Fold frequency in blinds versus steal attempts.", false),
		statsMetricDef(stats.MetricFoldBBToSteal, "Fold BB to Steal", "metric.fold_bb_to_steal.help", "Fold frequency from BB versus steal opens.", false),
		statsMetricDef(stats.MetricFoldSBToSteal, "Fold SB to Steal", "metric.fold_sb_to_steal.help", "Fold frequency from SB versus steal opens.", false),
		// C-bet and response
		statsMetricDef(stats.MetricFlopCBet, "Flop CBet", "metric.flop_cbet.help", "Continuation bet frequency on flop as preflop aggressor.", false),
		statsMetricDef(stats.MetricTurnCBet, "Turn CBet", "metric.turn_cbet.help", "Continuation bet frequency on turn.", false),
		statsMetricDef(stats.MetricDelayedCBet, "Delayed CBet", "metric.delayed_cbet.help", "Delayed continuation bet frequency (check flop, bet turn).", false),
		statsMetricDef(stats.MetricFoldToFlopCBet, "Fold to Flop CBet", "metric.fold_to_flop_cbet.help", "Fold frequency when facing flop c-bet.", false),
		statsMetricDef(stats.MetricFoldToTurnCBet, "Fold to Turn CBet", "metric.fold_to_turn_cbet.help", "Fold frequency when facing turn c-bet.", false),
		// Showdown and postflop profile
		statsMetricDef(stats.MetricWTSD, "WTSD", "metric.wtsd.help", "Went to showdown after seeing flop.", false),
		statsMetricDef(stats.MetricWSD, "W$SD", "metric.w_sd.help", "Won money at showdown.", true),
		statsMetricDef(stats.MetricWWSF, "WWSF", "metric.wwsf.help", "Won when saw flop.", false),
		statsMetricDef(stats.MetricAFq, "AFq", "metric.afq.help", "Aggression frequency: (bet+raise)/(actions postflop).", false),
		statsMetricDef(stats.MetricAF, "AF", "metric.af.help", "Aggression factor: (bet+raise)/call.", false),
		// Result profile
		statsMetricDef(stats.MetricWonWithoutSD, "Won without SD", "metric.won_without_sd.help", "Won hand without reaching showdown.", false),
		statsMetricDef(stats.MetricBBPer100, "bb/100", "metric.bb_per_100.help", "Net big blinds won per 100 hands.", false),
	}
}

func vRateOrZero(s *stats.Stats, id stats.MetricID) float64 {
	if s == nil {
		return 0
	}
	v, ok := s.Metric(id)
	if !ok {
		return 0
	}
	return v.Rate
}
