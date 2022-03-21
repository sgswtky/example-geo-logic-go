package geohash

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/mmcloughlin/geohash"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/sgswtky/sandbox/sampleData"
	"github.com/sgswtky/sandbox/util"
	"os"
	"sort"
)

type GeoHash struct {
	Hash string
	Lat  float64
	Lng  float64
	Name string
}

func (gh GeoHash) GetKey() *datastore.Key {
	return &datastore.Key{
		Kind: "GeoHash",
		Name: gh.Hash,
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
		gh := &GeoHash{Hash: geohash.Encode(v.Lat, v.Lng), Lat: v.Lat, Lng: v.Lng, Name: v.Name}
		if _, err := cli.Put(ctx, gh.GetKey(), gh); err != nil {
			panic(err)
		}
	}

	/**
	 * NOTE: 周囲500mの近傍する場合
	Geohashの桁数	南北の距離	東西の距離
	1				4,989,600.00m	4,050,000.00m
	2				623,700.00m	1,012,500.00m
	3				155,925.00m	126,562.50m
	4				19,490.62m	31,640.62m
	5				4,872.66m	3,955.08m
	6				609.08m	988.77m
	7				152.27m	123.60m
	8				19.03m	30.90m
	9				4.76m	3.86m
	10				0.59m	0.97m
	*/
	const FILTER_DISTANCE = 500
	// 6桁で取得、横3、縦3
	//           609.08m
	// 988.77m, (988.77m), 988.77m
	//           609.08m

	// NOTE: このあたりの処理を取得するメートルによってどう取得するかを変える必要がある
	loc := sampleData.GetThisLocation()
	// 東西(988.77 * 3) * 南北(609.08m * 3) の距離
	neighbors := geohash.Neighbors(string([]rune(geohash.Encode(loc.Lat, loc.Lng))[:6]))
	searchBlocks := neighbors
	for _, v := range neighbors {
		searchBlocks = append(searchBlocks, geohash.Neighbors(v)...)
	}

	type GeoHashDistance struct {
		GeoHash  GeoHash
		Distance float64
	}

	var result []GeoHashDistance
	for _, searchBlock := range util.ToUniqStrings(searchBlocks) {
		var tmp []GeoHash
		q := datastore.NewQuery(GeoHash{}.GetKey().Kind).
			Filter("Hash >= ", searchBlock).
			Filter("Hash <", searchBlock+"\ufffd")
		if _, err := cli.GetAll(ctx, q, &tmp); err != nil {
			panic(err)
		}
		// orbで直線距離を算出してそれごとに持っておく
		for _, v := range tmp {
			p1 := orb.Point{v.Lng, v.Lat}
			p2 := orb.Point{loc.Lng, loc.Lat}
			distance := geo.Distance(p1, p2)
			if FILTER_DISTANCE >= distance {
				result = append(result, GeoHashDistance{v, distance})
			}
		}
	}
	// NOTE: 距離順にして1件目を取ると一番近いものが取れる
	sort.Slice(result, func(i, j int) bool { return result[i].Distance < result[j].Distance })

	for _, v := range result {
		fmt.Printf("%fm, %s, %s\n", v.Distance, v.GeoHash.Hash, v.GeoHash.Name)
	}
}
