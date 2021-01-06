package main

import (
	"testing"
)

var (
	testDataPayloads = []struct {
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
)

func Test_Parse(t *testing.T) {
	for _, data := range testDataPayloads {
		cmd, err := NewRespReader([]byte(data.Payload))
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

func Benchmark_Parse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewRespReader([]byte("*4\r\n\r\n$6\r\nLRANGE\r\n$6\r\nmylist\r\n$1\r\n0\r\n$3\r\n599\r\n"))
	}
}
