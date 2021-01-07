package main

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"testing"
)

var (
	_ = []struct {
		Payload string
		Cmd     string
		Args    string
	}{
		{"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", "foo", "bar"},
		{"*2\r\n$4\r\nping\r\n$2\r\nok", "ping", "ok"},
		{"*1\r\n$4\r\nping\r\n", "ping", ""},
		{"*3\r\n$5\r\nRPUSH\r\n$6\r\nmylist\r\n$3\r\none\r\n", "RPUSH", "mylist one"},
		{"*4\r\n$6\r\nLRANGE\r\n$6\r\nmylist\r\n$1\r\n0\r\n$3\r\n599\r\n", "LRANGE", "mylist 0 599"},
		{"*4\r\n$6\r\nLRANGE\r\n$6\r\nmylist\r\n$1\r\n2a340d0a24340d0a485345540d0a2435340d0a75736572733\r\n$3\r\n2a340d0a24340d0a485345540d0a2435340d0a75736572733XXXXXXX\r\n", "LRANGE", "mylist 2a340d0a24340d0a485345540d0a2435340d0a75736572733 2a340d0a24340d0a485345540d0a2435340d0a75736572733"},
	}

	testDataPayloads = []struct {
		Payload string
		Cmd     string
		Args    string
		Sep     []byte
		Cls     string
		Size    int
	}{
		{"*3\r\n$5\r\nRPUSH\r\n$6\r\nuser:slot:ghost:095b314d-8e62-4e6c-abd6-e8a826ace563:chess\r\n$3\r\none\r\n",
			"RPUSH",
			"user:slot:ghost",
			[]byte{':'},
			`[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`,
			100,
		}, {"*2\r\n$3\r\nGET\r\n$56\r\nusernamechangecount_747ca354-ba9c-4c90-8e05-9a0dfe4ff668",
			"GET",
			"usernamechangecount_",
			[]byte{':'},
			`[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`,
			100},
	}
)

func Test_Parse(t *testing.T) {
	for _, data := range testDataPayloads {
		cls := regexp.MustCompile(data.Cls)
		cmd, err := NewRespReader([]byte(data.Payload), data.Sep, cls, data.Size)
		if err != nil {
			t.Error(err)
		}

		if cmd.Command() != data.Cmd {
			t.Errorf("cmd error. expect: %s, got: %s", data.Cmd, cmd.Command())
		}

		if cmd.Args() != data.Args {
			t.Errorf("args error. expect: %s, got: %s", data.Args, cmd.Args())
		}
	}
}

func Test_ParseFomatKey(t *testing.T) {

	testDatas := []struct {
		Payload string
		Pattern string
		Sep     string
		Expect  string
	}{
		{"user:slot:ghost", "slot", ":", "user::ghost"},
		{"user:slot:095b314d-8e62-4e6c-abd6-e8a826ace563:spagetti", `[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`, ":", "user:slot::spagetti"},
	}

	for _, data := range testDatas {
		rsp := RespReader{payload: []byte(data.Payload)}
		pattern := regexp.MustCompile(data.Pattern)
		out := rsp.cleanMatched(rsp.payload, pattern, []byte(data.Sep))
		outstr := string(out)
		if outstr != data.Expect {
			t.Errorf("clean error expect: %s, got: %s", data.Expect, outstr)
		}
	}
}

func Benchmark_ParseWithEmptyOptions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewRespReader([]byte("*4\r\n\r\n$6\r\nLRANGE\r\n$6\r\nmylist\r\n$1\r\n0\r\n$3\r\n599\r\n"), []byte{}, nil, -1)
	}
}

func Benchmark_ParseWithOptions(b *testing.B) {
	pattern := regexp.MustCompile(`\d...`)
	sep := []byte{':'}
	for i := 0; i < b.N; i++ {
		_, _ = NewRespReader([]byte("*4\r\n\r\n$6\r\nLRANGE\r\n$6\r\na:b:5000:dxxxxxxxxxxxxxxxxxxxx\r\n$1\r\n0\r\n$3\r\n599\r\n"),
			sep,
			pattern,
			20)
	}
}
