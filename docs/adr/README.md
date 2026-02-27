# Architecture Decision Records

このディレクトリ（`docs/adr/`）に ADR を管理します。

## OpenCode・エージェント向け

ADR の全体像（ステータス・日時）を把握するには以下を実行します。

    mise run adr-list

## 人間向けの操作

| 操作 | コマンド |
|---|---|
| 一覧表示 | `mise run adr-list` |
| 新規作成 | `adrgen create "決定タイトル"` |
| ステータス変更 | `adrgen status <ID> <status>` |
| 置き換え関係 | `adrgen create "新タイトル" -s <旧ID>` |
| 改訂関係 | `adrgen create "新タイトル" -a <旧ID>` |

## ステータス一覧

| 値 | 意味 |
|---|---|
| proposed | 提案中 |
| accepted | 採用・現行 |
| deprecated | 廃止済み |
| superseded | 別 ADR に置き換え済み |
| amended | 別 ADR により改訂済み |

## adrgen について

- ADR のステータス一覧を把握する場合は `mise run adr-list` を推奨
- `adrgen list` を直接使っても良いが、ツール導入は `mise install` が前提
