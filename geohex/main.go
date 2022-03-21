package geohex

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/bsm/go-geohex/v3"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/sgswtky/sandbox/sampleData"
	"os"
	"sort"
)

type GeoHex struct {
	Code string
	Lat  float64
	Lng  float64
	Name string
}

func (gh GeoHex) GetKey() *datastore.Key {
	return &datastore.Key{
		Kind: "GeoHex",
		Name: gh.Code,
	}
}

func Example() {
	_ = os.Setenv("DATASTORE_EMULATOR_HOST", "localhost:8081")

	ctx := context.Background()
	projectID := "test-qa"

	cli, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}

	for _, v := range sampleData.GetSampleData() {
		pos, _ := geohex.Encode(v.Lat, v.Lng, 20)
		gh := GeoHex{Code: pos.Code(), Lat: v.Lat, Lng: v.Lng, Name: v.Name}
		if _, err := cli.Put(ctx, gh.GetKey(), &gh); err != nil {
			panic(err)
		}
	}

	/**
	 * 半径 500mで取得する場合
	 *
	 * ズームレベル⇔直径or半径の対応表などは見当たらなかったので、
	 * GoogleMapsで直線距離を図って論理的に取得できているであろうロジックを組み立てて検証
	 *
	 * ズームレベル8で半径175mくらい取得できてそうなのでこれを基準にしています。
	 * （もしかしたら結構ズレあるかもなので、採用時には詳細な調査になりそうです）
	 */

	const FILTER_DISTANCE = 500

	// NOTE: このあたりの処理を取得するメートルによってどう取得するかを変える必要がある
	loc := sampleData.GetThisLocation()
	locPos, _ := geohex.Encode(loc.Lat, loc.Lng, 8)
	code := []rune(locPos.Code())
	// 親ブロックで検索すると 175*2=350 の半径をおおよそ検索可能
	locParent, _ := geohex.Decode(string(code[:len(code)-1]))
	// 親と親の周囲のブロックを検索する事で 350*2=700 の半径をおおよそ正確に検索可能
	searchBlocks := []string{locParent.Code()}
	for _, v := range locParent.Neighbours() {
		searchBlocks = append(searchBlocks, v.Code())
	}

	type GeoHexDistance struct {
		GeoHex   GeoHex
		Distance float64
	}

	var result []GeoHexDistance
	for _, searchBlock := range searchBlocks {
		query := datastore.NewQuery(GeoHex{}.GetKey().Kind).
			Filter("Code >= ", searchBlock).
			Filter("Code <", searchBlock+"\ufffd")
		var tmp []GeoHex
		if _, err := cli.GetAll(ctx, query, &tmp); err != nil {
			panic(err)
		}
		// orbで直線距離を算出してそれごとに持っておく
		for _, v := range tmp {
			p1 := orb.Point{v.Lng, v.Lat}
			p2 := orb.Point{loc.Lng, loc.Lat}
			distance := geo.Distance(p1, p2)
			// フィルタせずに確認すると 1000m超えるものも取得できているが、取得の範囲が円ではないので誤差で取得されてそう
			if FILTER_DISTANCE >= distance {
				result = append(result, GeoHexDistance{v, distance})
			}
		}
	}
	// NOTE: 距離順にして1件目を取ると一番近いものが取れる
	sort.Slice(result, func(i, j int) bool { return result[i].Distance < result[j].Distance })

	for _, v := range result {
		fmt.Printf("%fm, %s, %s\n", v.Distance, v.GeoHex.Code, v.GeoHex.Name)
	}
}
