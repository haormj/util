//go:build linux
// +build linux

package util

import (
	"net"
	"syscall"
)

const (
	SO_ORIGINAL_DST      = 80 // from linux/include/uapi/linux/netfilter_ipv4.h
	IP6T_SO_ORIGINAL_DST = 80 // from linux/include/uapi/linux/netfilter_ipv6/ip6_tables.h
)

// GetOriginalDst Get the original destination of a TCP connection.
// inspired by github.com/missdeer/avege
func GetOriginalDst(c *net.TCPConn) (string, uint16, error) {
	f, err := c.File()
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	fd := f.Fd()

	// The File() call above puts both the original socket fd and the file fd in blocking mode.
	// Set the file fd back to non-blocking mode and the original socket fd will become non-blocking as well.
	// Otherwise blocking I/O will waste OS threads.
	if err = syscall.SetNonblock(int(fd), true); err != nil {
		return "", 0, err
	}

	return getorigdst(fd)
}

// Call getorigdst() from linux/net/ipv4/netfilter/nf_conntrack_l3proto_ipv4.c
func getorigdst(fd uintptr) (string, uint16, error) {
	// IPv4 address starts at the 5th byte, 4 bytes long (206 190 36 45)
	addr, err := syscall.GetsockoptIPv6Mreq(int(fd), syscall.IPPROTO_IP, SO_ORIGINAL_DST)
	if err != nil {
		return "", 0, err
	}

	host := net.IP(addr.Multiaddr[4 : 4+net.IPv4len]).String()
	port := uint16(addr.Multiaddr[2])<<8 + uint16(addr.Multiaddr[3])

	return host, port, nil
}
