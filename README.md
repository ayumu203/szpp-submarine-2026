# SZPP - 潜水艦撃沈ゲーム開発

## 目的

- B1: この1年間で学んだHTML・JavaScriptの知識のおさらい・実践.
- B2: 	
	- ソフトウェアアーキテクチャの学習.
	- Go言語の考え方・チーム開発を抑える.
- B3: マネジメント・事故った箇所を巻き取る.

## 環境構築

```bash
git clone https://github.com/ayumu203/szpp-submarine-2026.git
```

### フロントエンド

- HTML + JavaScript + jQuery(jQuery Core 3.7.1 minified)

### バックエンド

- Go(1.22.2)
- 標準ライブラリのみで問題ないと思う.

```bash
sudo apt install gotestsum
```

### DB

- Upstash
- Redis を扱えるサーバレスのデータベース.


### システムのアーキテクチャ構成
```
| フロントエンド | ー(HTTP Request)→ |バックエンド| ー(HTTP Request)→ | Upstash |
```

![Tech Stack](https://skillicons.dev/icons?i=linux,ubuntu,html,js,jquery,go,redis,googlecloud)
