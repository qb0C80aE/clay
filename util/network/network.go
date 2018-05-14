package network

import (
	"encoding/binary"
	"fmt"
	"github.com/qb0C80aE/clay/logging"
	"net"
)

var utility = &Utility{}

// Utility handles network operation
type Utility struct {
}

// GetUtility returns the instance of utility
func GetUtility() *Utility {
	return utility
}

// ParseCIDR parses CIDR string and generate Ipv4Address instance based on that
func (receiver *Utility) ParseCIDR(cidr string) (*Ipv4Address, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		logging.Logger().Debug(err.Error())
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

// Ipv4Address is the struct represents IPv4 what is easy to use in calculation
type Ipv4Address struct {
	v4address uint32
	mask      uint32
	prefix    uint
	network   uint32
	host      uint32
}

func (receiver *Ipv4Address) clone() *Ipv4Address {
	result := &Ipv4Address{
		v4address: receiver.v4address,
		mask:      receiver.mask,
		prefix:    receiver.prefix,
		network:   receiver.network,
		host:      receiver.host,
	}
	return result
}

// IncreaseHostAddress increases given Ipv4Address host part considering its netmask
func (receiver *Ipv4Address) IncreaseHostAddress() *Ipv4Address {
	result := receiver.clone()
	if (result.host & ^result.mask) == (0xFFFFFFFF & ^result.mask) {
		result.host = 0
	} else {
		result.host++
	}
	return result
}

// DecreaseHostAddress decreases given Ipv4Address host part considering its netmask
func (receiver *Ipv4Address) DecreaseHostAddress() *Ipv4Address {
	result := receiver.clone()
	if result.host == 0x00000000 {
		result.host = 0xFFFFFFFF & ^result.mask
	} else {
		result.host--
	}
	return result
}

// IncreaseNetworkAddress increases given Ipv4Address network part
func (receiver *Ipv4Address) IncreaseNetworkAddress() *Ipv4Address {
	result := receiver.clone()
	result.network++
	return result
}

// DecreaseNetworkAddress decreases given Ipv4Address network part
func (receiver *Ipv4Address) DecreaseNetworkAddress() *Ipv4Address {
	result := receiver.clone()
	result.network--
	return result
}

// IncreaseIPAddress increases given Ipv4Address considering its network part and host part with netmask
func (receiver *Ipv4Address) IncreaseIPAddress() *Ipv4Address {
	result := receiver.clone()
	if (result.host & ^result.mask) == (0xFFFFFFFF & ^result.mask) {
		result.host = 0
		result.network++
	} else {
		result.host++
	}
	return result
}

// DecreaseIPAddress decreases given Ipv4Address considering its network part and host part with netmask
func (receiver *Ipv4Address) DecreaseIPAddress() *Ipv4Address {
	result := receiver.clone()
	if result.host == 0x00000000 {
		result.host = 0xFFFFFFFF & ^result.mask
		result.network--
	} else {
		result.host--
	}
	return result
}

// LimitedBroadcastAddress calculates the limited broadcast address from given Ipv4Address
func (receiver *Ipv4Address) LimitedBroadcastAddress() *Ipv4Address {
	result := receiver.clone()
	result.host = 0xFFFFFFFF & ^result.mask
	return result
}

// NetworkAddress calculates the network address from given Ipv4Address
func (receiver *Ipv4Address) NetworkAddress() *Ipv4Address {
	result := receiver.clone()
	result.host = 0
	return result
}

// MaxHostAddress calculates the maximum host address from given Ipv4Address
func (receiver *Ipv4Address) MaxHostAddress() *Ipv4Address {
	result := receiver.clone()
	result.host = (0xFFFFFFFF & ^result.mask) - 1
	return result
}

// MinimumHostAddress calculates the minimum host address from given Ipv4Address
func (receiver *Ipv4Address) MinimumHostAddress() *Ipv4Address {
	result := receiver.clone()
	result.host = 1
	return result
}

// IsBroadcastAddress checks if given Ipv4Address is the broadcast address or not
func (receiver *Ipv4Address) IsBroadcastAddress() bool {
	return ((receiver.host & ^receiver.mask) == (0xFFFFFFFF & ^receiver.mask))
}

// IsNetworkAddress checks if given Ipv4Address is network address or not
func (receiver *Ipv4Address) IsNetworkAddress() bool {
	return (receiver.host == 0x00000000)
}

// IsIncluding checks if left includes right network address
func (receiver *Ipv4Address) IsIncluding(target *Ipv4Address) bool {
	if target.prefix < receiver.prefix {
		return false
	}
	leftNetwork := receiver.network << uint(32-receiver.prefix)
	rightNetwork := target.network << uint(32-target.prefix)
	diff := rightNetwork & ^receiver.mask
	return ((leftNetwork | diff) == rightNetwork)
}

// String returns the string expression of given Ipv4Address
func (receiver *Ipv4Address) String() string {
	ip := (receiver.network << uint32(32-receiver.prefix)) | receiver.host
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, ip)
	return fmt.Sprintf("%d.%d.%d.%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
}

// NetMask returns the string expression of given Ipv4Address netmask
func (receiver *Ipv4Address) NetMask() string {
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, receiver.mask)
	return fmt.Sprintf("%d.%d.%d.%d", ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
}

// CIDR returns the string expression of given Ipv4Address with prefix
func (receiver *Ipv4Address) CIDR() string {
	return fmt.Sprintf("%s/%d", receiver.String(), receiver.prefix)
}
