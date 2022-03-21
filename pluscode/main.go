package pluscode

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	olc "github.com/google/open-location-code/go"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/sgswtky/sandbox/sampleData"
	"os"
	"sort"
)

type PlusCode struct {
	Code string
	Lat  float64
	Lng  float64
	Name string
}

func (pc PlusCode) GetKey() *datastore.Key {
	return &datastore.Key{
		Kind: "PlusCode",
		Name: pc.Code,
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

	geoCodes := sampleData.GetSampleData()

	for _, v := range geoCodes {
		pc := PlusCode{
			Code: olc.Encode(v.Lat, v.Lng, 20),
			Lat:  v.Lat,
			Lng:  v.Lng,
			Name: v.Name,
		}
		if _, err := cli.Put(ctx, pc.GetKey(), &pc); err != nil {
			panic(err)
		}
	}

	/**
	 * NOTE: 周囲500mの近傍する場合
	 * コード長	おおよそサイズ
	 *     2	2200km
	 *     4	110km
	 *     6	5.5km
	 *     8	275m
	 *     +
	 *     10	14m
	 *     11	3.5m
	 *
	 * 8桁で周囲2周分を取得する形で実装。
	 * ※TmpNeighbors で静的に返るようにして検証。採用する場合は開発する必要がある。
	 *  仕様がwikipediaで公開されてるので、これを確認しながら開発はできそう。
	 */
	const FILTER_DISTANCE = 500

	loc := sampleData.GetThisLocation()
	locPlusCode := olc.Encode(loc.Lat, loc.Lng, 8)
	searchBlocks := append([]string{locPlusCode}, TmpNeighbors(locPlusCode)...)

	type PlusCodeDistance struct {
		PlusCodeKind PlusCode
		Distance     float64
	}

	var result []PlusCodeDistance
	for _, searchBlock := range searchBlocks {
		query := datastore.NewQuery(PlusCode{}.GetKey().Kind).
			Order("Code").
			Filter("Code >= ", searchBlock).
			Filter("Code < ", searchBlock+"\ufffd")
		var tmp []PlusCode
		if _, err := cli.GetAll(ctx, query, &tmp); err != nil {
			panic(err)
		}
		for _, v := range tmp {
			p1 := orb.Point{v.Lng, v.Lat}
			p2 := orb.Point{loc.Lng, loc.Lat}
			distance := geo.Distance(p1, p2)
			if FILTER_DISTANCE >= distance {
				result = append(result, PlusCodeDistance{v, distance})
			}
		}
	}

	// NOTE: 距離順にして1件目を取ると一番近いものが取れる
	sort.Slice(result, func(i, j int) bool { return result[i].Distance < result[j].Distance })

	for _, v := range result {
		fmt.Printf("%fm, %s, %s\n", v.Distance, v.PlusCodeKind.Code, v.PlusCodeKind.Name)
	}
}

// TmpNeighbors NOTE: 周囲のブロックを取得するロジックを作成する必要がありそう
func TmpNeighbors(_ string) []string {
	return []string{
		"8Q5GH9VW+", // 左上
		"8Q5GH9VX+", // 上
		"8Q5GHCV2+", // 右上
		"8Q5GH9RW+", // 左
		"8Q5GHCR2+", // 右
		"8Q5GH9QW+", // 左下
		"8Q5GH9QX+", // 下
		"8Q5GHCQ2+", // 右下

		// 2週分使うのでここで返してしまう。本来はここに書かれてるべきではない
		"8Q5GH9WV+",
		"8Q5GH9WW+",
		"8Q5GH9WX+",
		"8Q5GHCW2+",
		"8Q5GHCW3+",
		"8Q5GHCV3+",
		"8Q5GHCR3+",
		"8Q5GHCQ3+",
		"8Q5GHCP3+",
		"8Q5GHCP2+",
		"8Q5GH9PX+",
		"8Q5GH9PW+",
		"8Q5GH9PV+",
		"8Q5GH9QV+",
		"8Q5GH9RV+",
		"8Q5GH9VV+",
	}
}
