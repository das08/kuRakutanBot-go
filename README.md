# kuRakutanBot-go
京都大学の講義別の単位取得率を検索できるLINEBotです

# 友だち追加
「京大楽単Bot」を友だち追加するには以下のボタンもしくはQRコードからお願いします。  

<img src="https://user-images.githubusercontent.com/41512077/101121375-590a2080-3633-11eb-9a10-bcdd3c4f7c2f.png" width="150px">  
<a href="http://nav.cx/eQHij4J"><img src="https://scdn.line-apps.com/n/line_add_friends/btn/ja.png" alt="友だち追加" height="36" border="0"></a>

# 使い方
- トーク画面で講義名を入力すると、前方一致でヒットする講義の単位取得率が返ってきます。  
<img src="https://user-images.githubusercontent.com/41512077/101122164-2c570880-3635-11eb-8d03-db913ebf4c51.jpg" width="150px">  

- 複数候補がある場合は目的の講義を選択することができます。  
<img src="https://user-images.githubusercontent.com/41512077/101122403-c919a600-3635-11eb-951b-7e33b67890c0.jpg" width="150px">  

- また、メニューからおみくじやアカデミックカレンダーが見れたりします。  
<img src="https://user-images.githubusercontent.com/41512077/101122162-2a8d4500-3635-11eb-8192-72c78e092347.jpg" width="200px">  

お気に入り機能もあるため、気になった講義はお気に入り保存することもできます。

# Environments
- docker 20.10.14
- docker compose 2.2.3

# Run
1. Clone this repository.
2. Build using docker-compose.
```shell
docker compose build
```
3. Run.
```shell
docker compose up
```

App runs on port 8081.
