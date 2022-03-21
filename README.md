
## 近傍検索ロジック確認 in Datastore

### ヘルプ

```bash.sh
# export GO111MODULE=on && go mod download
# go build && ./sandbox
Please slect : [pluscode geohash quadkey geohex]
```

### 実行

```bash.sh
# go build && ./sandbox geohash
102.536079m, wvuxp70by9s1, 福岡PARCO 新館
169.631315m, wvuxp73sv3z2, 福岡中央郵便局
223.725752m, wvuxp6cbmh8n, ローソン 西鉄福岡天神駅改札内店
230.344607m, wvuxp6gr03yy, ローソン 福岡市役所店
247.640457m, wvuxp725ft2h, 博多らーめん ShinShin 天神本店
303.256366m, wvuxp7e05vzq, ファミリーマート 福岡天神四丁目店
317.229813m, wvuxp7hj5n83, アクロス福岡
336.820000m, wvuxp7deyrxm, セブンイレブン 福岡天神４丁目店
354.399400m, wvuxp5nu67zz, セブンイレブン福岡舞鶴店
426.090338m, wvuxp62pvbpx, Apple福岡
448.079377m, wvuxp6s57weu, 天神南駅
485.327021m, wvuxp5zf7xbh, ローソン 天神北店
```

### 内容

Datastore に保存された `sampleData.GetSampleData()` に対して、500メートルの近傍検索を中心位置 `sampleData.GetThisLocation()` として取得する。

これを各ロジックで実装しています。

各コードは大体以下のような内容です。

1. 予めサンプルデータの緯度経度をアルゴリズムにより文字列化してDatastoreに保存
2. 取得する中心位置の周囲ブロック（1~2周）を半径500mに達するよう取得
3. ブロックに入り切るレコードをDatastoreから前方一致で取得
4. 緯度経度から直線距離を算出して、500m 以内のものを採用
5. 採用されたレコードを出力
