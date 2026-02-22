package persistence

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/AkatukiSora/vrc-vrpoker-ststs/internal/parser"
)

func TestSaveImportBatchParity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		newRepo func(t *testing.T) ImportRepository
	}{
		{
			name: "memory",
			newRepo: func(_ *testing.T) ImportRepository {
				return NewMemoryRepository()
			},
		},
		{
			name: "sqlite",
			newRepo: func(t *testing.T) ImportRepository {
				repo, err := NewSQLiteRepository(filepath.Join(t.TempDir(), "stats.db"))
				if err != nil {
					t.Fatalf("new sqlite repo: %v", err)
				}
				t.Cleanup(func() {
					_ = repo.Close()
				})
				return repo
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.newRepo(t)
			batchRepo, ok := repo.(ImportBatchRepository)
			if !ok {
				t.Fatalf("repo does not implement ImportBatchRepository")
			}

			hand := &parser.Hand{
				ID:         1,
				StartTime:  time.Date(2026, 2, 21, 0, 0, 0, 0, time.UTC),
				EndTime:    time.Date(2026, 2, 21, 0, 0, 5, 0, time.UTC),
				Players:    map[int]*parser.PlayerHandInfo{0: {SeatID: 0}},
				IsComplete: true,
			}
			source := HandSourceRef{
				SourcePath: "test.log",
				StartByte:  0,
				EndByte:    128,
				StartLine:  1,
				EndLine:    5,
			}
			source.HandUID = GenerateHandUID(hand, source)

			cursor := ImportCursor{
				SourcePath:     source.SourcePath,
				NextByteOffset: source.EndByte,
				NextLineNumber: source.EndLine,
				UpdatedAt:      time.Now(),
			}

			res, err := batchRepo.SaveImportBatch(context.Background(), []PersistedHand{{Hand: hand, Source: source}}, cursor)
			if err != nil {
				t.Fatalf("first save import batch: %v", err)
			}
			if res.Inserted != 1 || res.Updated != 0 {
				t.Fatalf("first upsert result: %+v", res)
			}

			res, err = batchRepo.SaveImportBatch(context.Background(), []PersistedHand{{Hand: hand, Source: source}}, cursor)
			if err != nil {
				t.Fatalf("second save import batch: %v", err)
			}
			if res.Updated != 1 {
				t.Fatalf("second upsert should update existing row: %+v", res)
			}

			saved, err := repo.GetCursor(context.Background(), source.SourcePath)
			if err != nil {
				t.Fatalf("get cursor: %v", err)
			}
			if saved == nil || saved.NextByteOffset != source.EndByte {
				t.Fatalf("cursor not saved correctly: %+v", saved)
			}
		})
	}
}
