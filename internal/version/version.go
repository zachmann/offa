package version

import (
	_ "embed" // for go:embed
	"strconv"
	"strings"
)

const (
	SOFTWAREID = "993a3850-60bf-11f0-ac62-5b2200a58e68"
)

// VERSION holds the server's version
//
//go:embed VERSION
var VERSION string

// Version segments
var (
	MAJOR int
	MINOR int
	FIX   int
	PRE   int
)

func init() {
	if VERSION[len(VERSION)-1] == '\n' {
		VERSION = VERSION[:len(VERSION)-1]
	}
	v := strings.Split(VERSION, ".")
	MAJOR, _ = strconv.Atoi(v[0])
	MINOR, _ = strconv.Atoi(v[1])
	ps := strings.Split(v[2], "-")
	FIX, _ = strconv.Atoi(ps[0])
	if len(ps) > 1 {
		pre := strings.TrimPrefix(ps[1], "pr")
		PRE, _ = strconv.Atoi(pre)
	}
}
