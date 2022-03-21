package main

import (
	"flag"
	"fmt"
	"github.com/sgswtky/sandbox/geohash"
	"github.com/sgswtky/sandbox/geohex"
	"github.com/sgswtky/sandbox/pluscode"
	"github.com/sgswtky/sandbox/quadkey"
	"os"
)

func main() {
	algorightmPluscode, algorightmGeohash, algorightmQuadkey, algorightmGeohex := "pluscode", "geohash", "quadkey", "geohex"
	algorightms := []string{algorightmPluscode, algorightmGeohash, algorightmQuadkey, algorightmGeohex}
	errmsg := fmt.Sprintf("Please slect : %v\n", algorightms)
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Print(errmsg)
		os.Exit(1)
	}

	switch flag.Args()[0] {
	case algorightmPluscode:
		pluscode.Example()
	case algorightmGeohash:
		geohash.Example()
	case algorightmQuadkey:
		quadkey.Example()
	case algorightmGeohex:
		geohex.Example()
	default:
		fmt.Print(errmsg)
		fmt.Println(flag.Args()[0])
		os.Exit(1)
	}
}
