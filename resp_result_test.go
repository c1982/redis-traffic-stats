package main

import (
	"testing"
)

var (
	testData = []struct {
		Cmd  string
		Data []struct {
			Args string
			Size float64
		}
		ExpectCount      uint64
		ExpectCommonKeys []string
	}{
		{
			"SET",
			[]struct {
				Args string
				Size float64
			}{
				{"user:1 username1", 1},
				{"user:2 username2", 1},
				{"user:3 username2", 1},
				{"user:4 username2", 1},
				{"users:boost1", 1},
				{"users:boost2", 1},
				{"users:boost3", 1},
				{"users:slot:ghost", 1},
			},
			8,
			[]string{"user:", "users:boost", "users:slot:ghost"},
		},
		{
			"GET",
			[]struct {
				Args string
				Size float64
			}{
				{"user:1", 1},
				{"user:2", 1},
				{"user:2", 1},
				{"user:2", 1},
			},
			4,
			[]string{"user:"},
		},
	}
)

func Test_RespResult(t *testing.T) {
	results := NewRespResult()

	for i := 0; i < len(testData); i++ {
		d := testData[i]
		for a := 0; a < len(d.Data); a++ {
			results.Add(d.Cmd, d.Data[a].Args, d.Data[a].Size)
		}
	}

	for i := 0; i < len(testData); i++ {
		d := testData[i]
		stats := results.calculateCommandStats(d.Cmd)
		if stats.Count != d.ExpectCount {
			t.Errorf("count error. expect: %d, got: %d", d.ExpectCount, stats.Count)
		}

		if len(d.ExpectCommonKeys) != len(stats.Arguments) {
			t.Errorf("common keys count error. expect: %d, got: %d", len(d.ExpectCommonKeys), len(stats.Arguments))
			t.Logf("common keys: %+v (%s)", stats.Arguments, d.Cmd)
		}

		for x := 0; x < len(d.ExpectCommonKeys); x++ {
			ekey := d.ExpectCommonKeys[x]
			ok := func(e string) bool {
				for y := 0; y < len(stats.Arguments); y++ {
					if e == stats.Arguments[y].Argument {
						return true
					}
				}
				return false
			}(ekey)

			if !ok {
				t.Errorf("expected key not found. expected: %s", ekey)
			}
		}
	}

}
