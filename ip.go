package main

import (
	"net"

	"inet.af/netaddr"
)

type Addr struct {
	addr  netaddr.IP
	inter string
}

func externalIPs() ([]Addr, error) {
	var ips []Addr
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return ips, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			nip, ok := netaddr.FromStdIP(ip)
			if !ok {
				continue
			}

			if nip.Is6() {
				continue
			}

			ips = append(ips, Addr{
				addr:  nip,
				inter: iface.Name,
			})
		}
	}
	return ips, nil
}
