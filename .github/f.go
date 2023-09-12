package main

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func main() {
	links, _ := netlink.LinkList()
	for _, link := range links {
		stats := link.Attrs().Statistics
		fmt.Printf("%s: %d bytes in, %d bytes out\n", link.Attrs().Name, stats.RxBytes, stats.TxBytes)
	}
}
