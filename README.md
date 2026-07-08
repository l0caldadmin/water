# water

`water` is a native Go library for [TUN/TAP](http://en.wikipedia.org/wiki/TUN/TAP) interfaces.

`water` is designed to be simple and efficient. It

* wraps almost only syscalls and uses only Go standard types;
* exposes standard interfaces; plays well with standard packages like `io`, `bufio`, etc.;
* does not handle memory management (allocating/destructing slice). It's up to user to decide whether/how to reuse buffers.

## Supported Platforms

| Platform | TUN | TAP |
|----------|-----|-----|
| Linux    | ✅  | ✅  |
| Windows  | ✅ (Wintun) | ❌ |
| macOS    | ✅ (point-to-point) | ❌ |
| FreeBSD / OpenBSD / NetBSD | ✅ | ✅ |

## Installation

```
go get github.com/l0caldadmin/water
```

## Documentation

[https://pkg.go.dev/github.com/l0caldadmin/water](https://pkg.go.dev/github.com/l0caldadmin/water)

## Examples

### TAP on Linux

```go
package main

import (
	"log"

	"github.com/l0caldadmin/water/waterutil"
	"github.com/l0caldadmin/water"
)

func main() {
	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = "O_O"

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}

	frame := make([]byte, 1500)
	for {
		n, err := ifce.Read(frame)
		if err != nil {
			log.Fatal(err)
		}
		buf := frame[:n]
		log.Printf("Dst: %s", waterutil.MACDestination(buf))
		log.Printf("Src: %s", waterutil.MACSource(buf))
	}
}
```

Bring it up:

```bash
sudo go run main.go &
sudo ip addr add 10.1.0.10/24 dev O_O
sudo ip link set dev O_O up
ping -c1 -b 10.1.0.255
```

### TUN on macOS

```go
package main

import (
	"log"

	"github.com/l0caldadmin/water"
)

func main() {
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	packet := make([]byte, 2000)
	for {
		n, err := ifce.Read(packet)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Packet Received: % x\n", packet[:n])
	}
}
```

```bash
sudo go run main.go &
sudo ifconfig utun2 10.1.0.10 10.1.0.20 up
ping 10.1.0.20
```

#### macOS caveats

1. Only point-to-point TUN devices are supported natively. TAP is not supported.
2. Custom interface names are not supported. Names are assigned automatically as `utun<N>`.

### TUN on Windows (Wintun)

Windows uses the [Wintun](https://www.wintun.net/) driver. TAP mode is not supported on this backend.

```go
package main

import (
	"log"

	"github.com/l0caldadmin/water"
)

func main() {
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	packet := make([]byte, 2000)
	for {
		n, err := ifce.Read(packet)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Packet Received: % x\n", packet[:n])
	}
}
```

Run as Administrator. Assign an IP address after the interface appears:

```dos
netsh interface ip set address name="wintun" source=static addr=10.1.0.10 mask=255.255.255.0 gateway=none
```

## Alternatives

* `tuntap`: [https://code.google.com/p/tuntap/](https://code.google.com/p/tuntap/)
