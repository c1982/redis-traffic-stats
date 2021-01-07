# redis-traffic-stats

### Installation

## Usage

>./redis-traffic-stats --interface=ens5 --addr=:9200 --password=pass --debug=false -s=: -r="[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}"

Parameters | Default | Desc
--- | --- | ---
debug | Enable |
interface | empty |
redisport | 6379 | 
addr | :9200 | 
username | admin | 
password | **** |

### Static Compilation

This tool require libpcap library (https://www.tcpdump.org/). 
You can embed dependencies on the compilation process. This helps for portability.
Check the build command below

> LDFLAGS='-l/usr/lib/libpcap.a' CGO_ENABLED=1 go build -ldflags '-linkmode external -extldflags -static' .

### Troubleshooting

If you get an error about libpcap like below

```bash
./redis-traffic-stats: error while loading shared libraries: libpcap.so.0.8: cannot open shared object file: No such file or directory
```

install libpcap

for centos
>yum install libpcap-devel

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

Oğuzhan YILMAZ - [@c1982](https://twitter.com/c1982) - aspsrc@gmail.com