package parser

import (
	"fmt"
	"os"
	"testing"
)

func TestDebugStats(t *testing.T) {
	logPath := "/home/akatuki-sora/.local/share/Steam/steamapps/compatdata/438100/pfx/drive_c/users/steamuser/AppData/LocalLow/VRChat/VRChat/output_log_2026-02-20_22-38-17.txt"
	f, err := os.Open(logPath)
	if err != nil {
		t.Skipf("log file not found: %v", err)
	}
	defer f.Close()

	result, _ := ParseReader(f)
	localSeat := result.LocalPlayerSeat
	t.Logf("=== Local Seat: %d | Total Hands: %d ===", localSeat, len(result.Hands))

	vpip, pfr, threebet, foldTo3Bet := 0, 0, 0, 0
	threeBetOpp, foldTo3BetOpp := 0, 0
	won, showdown, wonSD := 0, 0, 0
	totalPotWon, totalInvested := 0, 0

	for _, h := range result.Hands {
		pi, ok := h.Players[localSeat]
		if !ok {
			continue
		}

		// 投資額計算
		invested := 0
		for _, act := range pi.Actions {
			invested += act.Amount
		}
		totalInvested += invested
		totalPotWon += pi.PotWon

		if pi.Won {
			won++
		}
		if pi.ShowedDown {
			showdown++
			if pi.Won {
				wonSD++
			}
		}
		if pi.VPIP {
			vpip++
		}
		if pi.PFR {
			pfr++
		}
		if pi.ThreeBet {
			threebet++
			threeBetOpp++
		}
		if pi.FoldTo3Bet {
			foldTo3Bet++
			foldTo3BetOpp++
		}
		// 3bet機会: PFRした後に相手が3betした場合 or 3betした場合
		// → calculatePreflopStatsの内部ロジックを確認

		// 詳細ログ (最初の20手)
		if h.ID <= 20 {
			cards := ""
			if len(pi.HoleCards) == 2 {
				cards = fmt.Sprintf("%s%s+%s%s", pi.HoleCards[0].Rank, pi.HoleCards[0].Suit, pi.HoleCards[1].Rank, pi.HoleCards[1].Suit)
			}
			board := ""
			for _, c := range h.CommunityCards {
				board += c.Rank + c.Suit + " "
			}
			result := "lost"
			if pi.Won {
				result = fmt.Sprintf("WON %d", pi.PotWon)
			}
			t.Logf("Hand#%02d [%s] pos=%-4s cards=%-8s board=%-18s %s | VPIP=%v PFR=%v 3bet=%v F3bet=%v SD=%v invested=%d",
				h.ID, h.StartTime.Format("15:04"), pi.Position, cards, board, result,
				pi.VPIP, pi.PFR, pi.ThreeBet, pi.FoldTo3Bet, pi.ShowedDown, invested)

			// プリフロップアクション詳細
			pfStr := ""
			for _, act := range pi.Actions {
				if act.Street == StreetPreFlop {
					pfStr += fmt.Sprintf("%s(%d) ", act.Action, act.Amount)
				}
			}
			t.Logf("  PF actions: %s | SBseat=%d BBseat=%d players=%d", pfStr, h.SBSeat, h.BBSeat, h.NumPlayers)
		}
	}

	total := len(result.Hands)
	t.Logf("=== Summary ===")
	t.Logf("Win rate (won/total): %d/%d = %.1f%%", won, total, pct(won, total))
	t.Logf("VPIP:       %d/%d = %.1f%%", vpip, total, pct(vpip, total))
	t.Logf("PFR:        %d/%d = %.1f%%", pfr, total, pct(pfr, total))
	t.Logf("3Bet:       %d/%d opp = %.1f%%", threebet, threeBetOpp, pct(threebet, threeBetOpp))
	t.Logf("FoldTo3Bet: %d/%d opp = %.1f%%", foldTo3Bet, foldTo3BetOpp, pct(foldTo3Bet, foldTo3BetOpp))
	t.Logf("W$SD:       %d/%d = %.1f%%", wonSD, showdown, pct(wonSD, showdown))
	t.Logf("TotalPotWon:  %d", totalPotWon)
	t.Logf("TotalInvested:%d", totalInvested)
	t.Logf("Net profit:   %+d", totalPotWon-totalInvested)
}

func pct(n, d int) float64 {
	if d == 0 {
		return 0
	}
	return float64(n) / float64(d) * 100
}

func TestDebugHandRange(t *testing.T) {
	logPath := "/home/akatuki-sora/.local/share/Steam/steamapps/compatdata/438100/pfx/drive_c/users/steamuser/AppData/LocalLow/VRChat/VRChat/output_log_2026-02-20_22-38-17.txt"
	f, err := os.Open(logPath)
	if err != nil {
		t.Skipf("log file not found: %v", err)
	}
	defer f.Close()

	result, _ := ParseReader(f)
	localSeat := result.LocalPlayerSeat

	t.Logf("=== Hand Range Data (all dealt hands) ===")
	for _, h := range result.Hands {
		pi, ok := h.Players[localSeat]
		if !ok || len(pi.HoleCards) != 2 {
			continue
		}
		c1, c2 := pi.HoleCards[0], pi.HoleCards[1]
		suited := c1.Suit == c2.Suit
		isPair := c1.Rank == c2.Rank
		comboType := "o"
		if isPair {
			comboType = ""
		} else if suited {
			comboType = "s"
		}
		action := "fold"
		if pi.VPIP && pi.PFR {
			action = "raise"
		} else if pi.VPIP {
			action = "call"
		}
		t.Logf("  Hand#%02d %s%s%s%s%-2s VPIP=%v PFR=%v FoldPF=%v pos=%s",
			h.ID, c1.Rank, c1.Suit, c2.Rank, c2.Suit, comboType,
			pi.VPIP, pi.PFR, pi.FoldedPF, pi.Position)
		_ = action
	}
}
