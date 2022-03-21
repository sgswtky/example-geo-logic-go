package h3

import (
	"fmt"
	"github.com/uber/h3-go"
)

func Example() {
	geo := h3.GeoCoord{
		Latitude:  33.593312,
		Longitude: 130.385187,
	}
	// TODO: resolution とは何か。これによって数字が全然変わる。
	//       15: 0x8f30d9ac6d2c826, 6: 0x8630d9af7ffffff
	//       上記のようになるため前方一致での探索は不可能。
	//resolution := 15
	//fmt.Printf("%#x\n", h3.FromGeo(geo, resolution))
	//fmt.Printf("%#x\n", h3.FromGeo(geo, 6))
	for i := 0; i < 15; i++ {
		fmt.Printf("%#x\n", h3.FromGeo(geo, i+1))
	}
	// TODO: 文字列削れば上の階層になるかどうかがまだわかってない。たどり着いてない。
	//       16進数で表されてるように見えるけど `f` で埋めたり `0` で埋めたりしても上位の階層を再現できるアルゴリズムではなさそう
	//       https://h3geo.org/docs/core-library/h3Indexing#bit-layout-of-h3index
}
