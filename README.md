# チャットアプリ
## 概要
- `WebSocket` の接続確認
```
go run main.go

curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" -H "Sec-WebSocket-Version: 13" -H "Sec-WebSocket-Key: WIY4slX50bnnSF1GaedKhg==" -H "Host: localhost:8080" http://localhost:8080/ws
```


## メモ
- `goroutine` や `channel` の実装練習
- 双方向通信プロトコルの `WebSocket` についてのキャッチアップ
- チャットアプリのようなユーザーA、ユーザーBが同時にサーバーに接続している状況において、リアルタイムにメッセージをやり取りするには、双方向通信のプロトコルである `WebSocket` が使える
- `WebSocket` は `ユーザーA - サーバー` , `ユーザーB - サーバー` で個別に接続を確立する
  - サーバーを介してユーザーAとユーザーBはリアルタイムにメッセージのやり取りを行う
  - ユーザーAが接続を切った後で、ユーザーBがメッセージを送信しユーザーBが接続を切り、その後再度ユーザーAが接続を確立するとサーバーから未取得のメッセージを取得する
  - `LINE` はアプリを起動していなくてもプッシュ通知でメッセージが来たことがわかるが、これは接続が常に確立されているわけではなく、iOSやAndroidのプッシュ通知サービスを利用しているだけで、ちゃんとメッセージを取得するにはアプリを立ち上げて接続を確立することが必要(とのこと by ChatGPT)
- `双方向通信プロトコル` は、従来のHTTP通信のリクエスト-レスポンスのような関係ではなく、クライアントからサーバー、サーバーからクライアントにいつでも情報を送ることができる
- [gorilla/websocket](https://github.com/gorilla/websocket) を使ってみる
