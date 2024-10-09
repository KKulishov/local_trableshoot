package net

import (
	"local_trableshoot/internal/format"
	"os"
)

func GetNetworkStats(file *os.File) {
	// Interfaces
	format.WriteHeader(file, "Interfaces")
	interfaceOutput := format.ExecuteCommand("ip", "-s", "link", "show")
	format.WritePreformatted(file, interfaceOutput)

	// Addresses
	format.WriteHeader(file, "Addresses")
	addressesOutput := format.ExecuteCommand("ip", "-s", "address", "show")
	format.WritePreformatted(file, addressesOutput)

	// Routes
	format.WriteHeader(file, "Routes")
	routesOutput := format.ExecuteCommand("ip", "-s", "route", "show")
	format.WritePreformatted(file, routesOutput)

	// IPv6 Routes (optional)
	haveIPv6 := os.Getenv("HAVE_IPV6")
	if haveIPv6 != "" {
		ipv6RoutesOutput := format.ExecuteCommand("ip", "-s", "-6", "route", "show")
		format.WritePreformatted(file, ipv6RoutesOutput)
	}

	// Neighbours
	format.WriteHeader(file, "Neighbours")
	neighboursOutput := format.ExecuteCommand("ip", "-s", "neighbour", "show")
	format.WritePreformatted(file, neighboursOutput)

	// IPv6 Neighbours (optional)
	if haveIPv6 != "" {
		ipv6NeighboursOutput := format.ExecuteCommand("ip", "-s", "-6", "neighbour", "show")
		format.WritePreformatted(file, ipv6NeighboursOutput)
	}

	// Resolvers
	format.WriteHeader(file, "Resolvers")
	resolversOutput := format.ExecuteCommand("cat", "/etc/resolv.conf")
	format.WritePreformatted(file, resolversOutput)
}
