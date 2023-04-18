# RssNews
RSSを登録して購読できるweb app (beta ver)

# name

トップページ
![image](./top.png)

メニュー
![image](./menu.png)

## Overview
golang/HTML/CSS/javascript/postgresqlを勉強し初めて作成したwebページです.

## Requirement

### langauge
japanese

### os
windows

### database
postgresql

### programming language
golang 
javascript

HTML/CSS

## Usage
動作にはpostgresqlが必要となります.
1.postgresqlにrss_reader_webというデータベースを作成する.
2.main.goを起動

で基本的にシステムは動作します.

## TODO
1. cookieのuserIDから登録したRSSURLなどをテーブルから入手しているのでこれを別の方法にする
2. ヘッダーとフッターなどの共通の部分のファイル切り分け
3. 現在タブのfeedに全ての記事を入れているのでサイトのタブに変更する
4. aboutページの実装
5. お気に入りボタンの実装(テーブルは実装済み)
6. favi iconの実装
7. sarinaを使用してlogoの実装(タイトルが思い浮かばないので)
8. 3が終了後に検索機能を実装
9. ホットリロードの実装

## Features
NULL

## Reference
NULL

## Author
NULL

