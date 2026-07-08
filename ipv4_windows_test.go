package water

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"testing"
)

// checkWindowsPrereqs skips the test when required runtime drivers are absent.
// TAP is not supported on this Windows backend (Wintun TUN only).
// TUN tests require wintun.dll to be loadable.
func checkWindowsPrereqs(t *testing.T, deviceType DeviceType) {
	t.Helper()
	if deviceType == TAP {
		t.Skip("TAP is not supported on Windows (Wintun TUN only)")
	}
	// Probe TUN availability without leaving a device open.
	ifce, err := New(Config{DeviceType: TUN})
	if err != nil {
		t.Skipf("skipping: Windows TUN prerequisite unavailable: %v", err)
	}
	ifce.Close()
}

func startPing(t *testing.T, dst net.IP, _ bool) {
	if err := exec.Command("ping", "-n", "4", dst.String()).Start(); err != nil {
		t.Fatal(err)
	}
}

func setupIfce(t *testing.T, ipNet net.IPNet, dev string) {
	args := strings.Fields("interface ip set address name=REPLACE_ME source=static addr=REPLACE_ME mask=REPLACE_ME gateway=none")
	args[4] = fmt.Sprintf("name=%s", dev)
	args[6] = fmt.Sprintf("addr=%s", ipNet.IP)
	args[7] = fmt.Sprintf("mask=%d.%d.%d.%d", ipNet.Mask[0], ipNet.Mask[1], ipNet.Mask[2], ipNet.Mask[3])
	cmd := exec.Command("netsh", args...)
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
}

func teardownIfce(t *testing.T, ifce *Interface) {
	if ifce == nil {
		return
	}
	if err := ifce.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestBroadcastTAP(t *testing.T) {
	checkWindowsPrereqs(t, TAP)

	var (
		self = net.IPv4(10, 0, 42, 1)
		mask = net.IPv4Mask(255, 255, 255, 0)
		brd  = net.IPv4(10, 0, 42, 255)
	)

	ifce, err := New(Config{DeviceType: TAP})
	if err != nil {
		t.Fatalf("creating TAP error: %v\n", err)
	}
	defer teardownIfce(t, ifce)

	dataCh, errCh := startRead(t, ifce)

	setupIfce(t, net.IPNet{IP: self, Mask: mask}, ifce.Name())
	startPing(t, brd, true)

	waitForPingOrBust(t, true, true, self, brd, dataCh, errCh)
}
