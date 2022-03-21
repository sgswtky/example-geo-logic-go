package quadkey

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/maptile"
	"github.com/sgswtky/sandbox/sampleData"
	"os"
	"sort"
	"strconv"
)

type QuadKey struct {
	QuadKey string
	Lat     float64
	Lng     float64
	Name    string
}

func (qkl QuadKey) GetKey() *datastore.Key {
	return &datastore.Key{
		Kind: "QuadKey",
		Name: qkl.QuadKey,
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
		i := maptile.At(orb.Point{v.Lng, v.Lat}, 28).Quadkey()
		qkl := QuadKey{
			// 4進数でString化して扱う
			QuadKey: strconv.FormatUint(i, 4),
			Lat:     v.Lat,
			Lng:     v.Lng,
			Name:    v.Name,
		}
		if _, err := cli.Put(ctx, qkl.GetKey(), &qkl); err != nil {
			panic(err)
		}
	}
	/**
	 * NOTE: 周囲500mの近傍する場合
	 * ズーム レベル	メートル/ピクセル	メートル/タイル一辺
	 * 0	156543	40075017
	 * 1	78271.5	20037508
	 * 2	39135.8	10018754
	 * 3	19567.88	5009377.1
	 * 4	9783.94	2504688.5
	 * 5	4891.97	1252344.3
	 * 6	2445.98	626172.1
	 * 7	1222.99	313086.1
	 * 8	611.5	156543
	 * 9	305.75	78271.5
	 * 10	152.87	39135.8
	 * 11	76.44	19567.9
	 * 12	38.219	9783.94
	 * 13	19.109	4891.97
	 * 14	9.555	2445.98
	 * 15	4.777	1222.99
	 * 16	2.3887	611.496
	 * 17	1.1943	305.748
	 * 18	0.5972	152.874
	 * 19	0.14929	76.437
	 * 20	0.14929	38.2185
	 * 21	0.074646	19.10926
	 * 22	0.037323	9.55463
	 * 23	0.0186615	4.777315
	 * 24	0.00933075	2.3886575
	 *
	 * ズームレベル16で1辺が 611mなのでコレで周囲を取得すれば 500mは取得可能
	 */
	const FILTER_DISTANCE = 500

	// NOTE: このあたりの処理を取得するメートルによってどう取得するかを変える必要がある
	loc := sampleData.GetThisLocation()
	locMapTile := maptile.At(orb.Point{loc.Lng, loc.Lat}, 16)
	searchTiles := []maptile.Tile{locMapTile}
	searchTiles = append(searchTiles, circumferenceNeighbors(locMapTile)...)

	type QuadKeyDistance struct {
		QuadKey  QuadKey
		Distance float64
	}

	var result []QuadKeyDistance
	for _, searchTile := range searchTiles {
		searchQuadKey := strconv.FormatUint(searchTile.Quadkey(), 4)
		var tmp []QuadKey
		q := datastore.NewQuery(QuadKey{}.GetKey().Kind).
			Filter("QuadKey >=", searchQuadKey).
			Filter("QuadKey <", searchQuadKey+"\ufffd")
		if _, err := cli.GetAll(ctx, q, &tmp); err != nil {
			panic(err)
		}
		// orbで直線距離を算出してそれごとに持っておく
		for _, v := range tmp {
			p1 := orb.Point{v.Lng, v.Lat}
			p2 := orb.Point{loc.Lng, loc.Lat}
			distance := geo.Distance(p1, p2)
			if FILTER_DISTANCE >= distance {
				result = append(result, QuadKeyDistance{v, distance})
			}
		}
	}

	// NOTE: 距離順にして1件目を取ると一番近いものが取れる
	sort.Slice(result, func(i, j int) bool { return result[i].Distance < result[j].Distance })

	for _, v := range result {
		fmt.Printf("%fm, %s, %s\n", v.Distance, v.QuadKey.QuadKey, v.QuadKey.Name)
	}
}

// circumferenceNeighbors ブロックを取り囲む縦横斜め8方向のブロックを取得
func circumferenceNeighbors(center maptile.Tile) maptile.Tiles {
	upper := UpperTile(center)
	left := LeftTile(center)
	right := RightTile(center)
	lower := LowerTile(center)
	return []maptile.Tile{
		// upper side
		LeftTile(upper),
		upper,
		RightTile(upper),

		left,
		right,

		// lower side
		LeftTile(lower),
		lower,
		RightTile(lower),
	}
}

func UpperTile(tile maptile.Tile) maptile.Tile {
	return maptile.New(tile.X, tile.Y-1, tile.Z)
}

func LeftTile(tile maptile.Tile) maptile.Tile {
	return maptile.New(tile.X-1, tile.Y, tile.Z)
}

func RightTile(tile maptile.Tile) maptile.Tile {
	return maptile.New(tile.X+1, tile.Y, tile.Z)
}

func LowerTile(tile maptile.Tile) maptile.Tile {
	return maptile.New(tile.X, tile.Y+1, tile.Z)
}
