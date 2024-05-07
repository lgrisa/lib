package netutil

import (
	"log"
	"net"
)

// Get preferred outbound ip of this machine
func GetOutboundIP(targetAddr string) net.IP {
	conn, err := net.Dial("udp", targetAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
