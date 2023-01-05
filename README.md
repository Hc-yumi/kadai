# 課題　 --

## ①課題内容（どんな作品か）
- ブックマークアプリ

## ②工夫した点・こだわった点
- Golangを使用してDB保存を実装した

## ④動作手順
### Postgresにログイン(要インストール)
データベースの接続コマンド。passwordは初期設定の値。
```
psql -h 127.0.0.1 -U postgres -W -p 5432
```
### データベースの作成
データベースの作成コマンド
```
postgres=# create database test;
```
main.goのコード中にある33行目のパスワードはpostgresインストール時に設定したもの
```
create table booklist(
id SERIAL PRIMARY KEY,
bookname varchar(64),
url varchar(100),
comment varchar(100),
time timestamp default now()
);
```

### GOの動かし方
GOの最新版(1.19)をインストールする。
そののちに以下のコマンドを実行
```
go mod init
go mod tidy
go run main.go
```

## ⑤質問・疑問・感想、シェアしたいこと等なんでも
- [質問]
- [疑問]
- [感想]GOを使用して最低限の課題を作成した。次回はもう少し機能を追加したい。
- [tips]
- [参考記事]
