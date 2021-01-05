# redis-traffic-stats

./redis-traffic-stats --interface=ens5 

./redis-traffic-stats: error while loading shared libraries: libpcap.so.0.8: cannot open shared object file: No such file or directory
> LDFLAGS='-l/usr/lib/libpcap.a' CGO_ENABLED=1 go build -ldflags '-linkmode external -extldflags -static' .

yum install libpcap-devel
yum install libpcap*
install libpcap0.8-dev libuv1-dev

```golang
duratios.m.Lock()
if tcp.SrcPort == redisPort {
	duration, ok := duratios.list[tcp.Seq]
	if ok {
		current := time.Now().UnixNano()
		latency := current - duration
		fmt.Printf("seq: %d, latency %s  len: %d\n", tcp.Seq, time.Nanosecond*time.Duration(latency), len(duratios.list))
		delete(duratios.list, tcp.Seq)
	}
}
duratios.list[tcp.Ack] = time.Now().UnixNano()
duratios.m.Unlock()
```