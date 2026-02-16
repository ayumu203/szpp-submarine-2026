# Git に関して

## 概要

- チーム開発をする上でT○○msやDisc○dでコードを共有するのはかなり開発の難易度を高くしてしまう(LINEは論外死刑).
- そんなときに`Git`, `GitHub`というツールが非常に役に立つ.

## 考え方

- Git ではそれぞれの作業ごとにブランチと呼ばれる枝をもとの流れから分岐させる.
- 作業が関連した枝を本流へ戻すことでチームで作業を分担して行うことができる.

## 基本

- リポジトリのクローン: インターネット上で共有されたプロジェクトを自分の環境へ複製する.

```bash
# git clone <url>
git clone https://github.com/szpp-dev-team/submarine-war-visualizer.git
```

- ログの確認: 今までの作業内容を見る(<font color="red">重要</font>).

```bash
# git log
git log
# こんなのが出てくると思う.
commit 70fdffb9b2c427037b3eba64b08a7c42dc15e2be (HEAD -> master)
Author: Ayumu <ochiohita@gmail.com>
Date:   Mon Feb 16 00:46:04 2026 +0900

    docs: 可視化へのリンクの追加.
・・・
```

- 状態の確認

```bash
# git status
git status
ブランチ master
追跡されていないファイル:
  (use "git add <file>..." to include in what will be committed)
	"materials/03_Git\345\210\251\347\224\250\346\226\271\346\263\225\343\201\250\343\203\253\343\203\274\343\203\253.md"

nothing added to commit but untracked files present (use "git add" to track)
```

- ブランチの移動(<font color="red">超重要</font>)

```bash
# git switch -c <ブランチ名> (ブランチを作成して移動)
git switch -c docs/create-readme
# git switch <ブランチ名> (ブランチを移動)
git switch dev
```

- 作業内容の保存(<font color="red">超重要</font>)

```bash
# git add ファイル名(変更の追加)
git add foo.md
# git commit -m "コミットメッセージ"(変更のコミット)
git commit -m "docs: プレイヤークラスを実装."
```

- リモート(本流)への変更の反映

```bash
git push
```
