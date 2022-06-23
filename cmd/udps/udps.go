package main

import (
	"flag"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "bind ip")
	port := flag.Int("port", 9797, "bind port")
	flag.Parse()

	IPAddress := net.ParseIP(*ip)
	ServerConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: IPAddress, Port: *port, Zone: ""})
	if err != nil {
		log.Errorf("can not create conn %s", err.Error())
		return
	}
	defer ServerConn.Close()
	buf := make([]byte, 1024)
	for {
		n, addr, _ := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)
		ServerConn.WriteToUDP(buf[0:n], addr)
	}
}
