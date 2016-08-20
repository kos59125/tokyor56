# このリポジトリについて

「第56回R勉強会＠東京（#TokyoR）」で発表した内容のソースコードです。

## ディレクトリ構成

```
tokyor56/
+ producer/:
+ consumer/:
```

## アプリケーション構成

### producer

アクセスログを生成するアプリケーションです。

```
# ビルド
cd producer/
glide install
go build

# 実行例
./producer -p 10 -o /var/log/producer/access.log
```

Kinesis へのデータ送信は本プログラムではサポートしていません。
任意の方法でログを Kinesis に送信する必要があります。

（例） IAM ロールを付与した EC2 インスタンスでプログラムを実行してログを生成し、 [Amazon Kinesis エージェント](http://docs.aws.amazon.com/ja_jp/streams/latest/dev/writing-with-agents.html)を利用して Kinesis にログを転送する

### consumer

Kinesis からデータを取得する shiny アプリケーションです。

