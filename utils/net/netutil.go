package netutil

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Ipv4Address struct {
	v4address uint32
	mask      uint32
	prefix    uint
	network   uint32
	host      uint32
}

func ParseCIDR(cidr string) (*Ipv4Address, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	ip = ip.To4()
	result := &Ipv4Address{}
	result.v4address = (uint32(ip[0])<<24 | (uint32(ip[1]) << 16) | (uint32(ip[2]) << 8) | uint32(ip[3]))
	prefix, _ := ipNet.Mask.Size()
	result.prefix = uint(prefix)
	result.mask = (uint32(ipNet.Mask[0])<<24 | (uint32(ipNet.Mask[1]) << 16) | (uint32(ipNet.Mask[2]) << 8) | uint32(ipNet.Mask[3]))
	result.network = (result.v4address & result.mask) >> uint32(32-prefix)
	result.host = (result.v4address & ^result.mask)
	return result, nil
}

func clone(source *Ipv4Address) *Ipv4Address {
	result := &Ipv4Address{
		v4address: source.v4address,
		mask:      source.mask,
		prefix:    source.prefix,
		network:   source.network,
		host:      source.host,
	}
	return result
}

func IncreaseHostAddress(source *Ipv4Address) *Ipv4Address {
	result := clone(source)
	if (result.host & ^result.mask) == (0xFFFFFFFF & ^result.mask) {
		result.host = 0
	} else {
		result.host += 1
	}
	return result
}

func DecreaseHostAddress(source *Ipv4Address) *Ipv4Address {
	result := clone(source)
	if result.host == 0x00000000 {
		result.host = 0xFFFFFFFF & ^result.mask
	} else {
		result.host -= 1
	}
	return result
}

func IncreaseNetworkAddress(source *Ipv4Address) *Ipv4Address {
	result := clone(source)
	result.network += 1
	return result
}

func DecreaseNetworkAddress(source *Ipv4Address) *Ipv4Address {
	result := clone(source)
	result.network -= 1
	return result
}

func IncreaseIpAddress(source *Ipv4Address) *Ipv4Address {
	result := clone(source)
	if (result.host & ^result.mask) == (0xFFFFFFFF & ^result.mask) {
		result.host = 0
		result.network += 1
	} else {
		result.host += 1
	}
	return result
}

func DecreaseIpAddress(source *Ipv4Address) *Ipv4Address {
	result := clone(source)
	if result.host == 0x00000000 {
		result.host = 0xFFFFFFFF & ^result.mask
		result.network -= 1
	} else {
		result.host -= 1
	}
	return result
}

func LimitedBroadcast(source *Ipv4Address) *Ipv4Address {
	result := clone(source)
	result.host = 0xFFFFFFFF & ^result.mask
	return result
}

func Network(source *Ipv4Address) *Ipv4Address {
	result := clone(source)
	result.host = 0
	return result
}

func MaxHost(source *Ipv4Address) *Ipv4Address {
	result := clone(source)
	result.host = (0xFFFFFFFF & ^result.mask) - 1
	return result
}

func MinimumHost(source *Ipv4Address) *Ipv4Address {
	result := clone(source)
	result.host = 1
	return result
}

func IsBroadcastAddress(source *Ipv4Address) bool {
	return ((source.host & ^source.mask) == (0xFFFFFFFF & ^source.mask))
}

func IsNetworkAddress(source *Ipv4Address) bool {
	return (source.host == 0x00000000)
}

func String(source *Ipv4Address) string {
	ip := (source.network << uint32(32-source.prefix)) | source.host
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, ip)
	return fmt.Sprintf("%d.%d.%d.%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
}

func StringWithPrefix(source *Ipv4Address) string {
	return fmt.Sprintf("%s/%d", String(source), source.prefix)
}
