package main

import (
	"flag"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var mu sync.Mutex

var sendTime, recvTime time.Time

func ConnRecv(conn *net.UDPConn) {
	buf := make([]byte, 1024)
	for {
		n, _, _ := conn.ReadFromUDP(buf)
		if n != 0 {
			mu.Lock()
			recvTime = time.Now()
			mu.Unlock()
		}
	}
}

func main() {
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	ip := flag.String("ip", "0.0.0.0", "bind ip")
	port := flag.Int("port", 7979, "bind port")

	rip := flag.String("rip", "0.0.0.0", "remote ip")
	rport := flag.Int("rport", 9797, "remote port")

	flag.Parse()

	IPAddress := net.ParseIP(*ip)
	RIPAddress := net.ParseIP(*rip)
	Conn, err := net.DialUDP("udp", &net.UDPAddr{IP: IPAddress, Port: *port, Zone: ""}, &net.UDPAddr{IP: RIPAddress, Port: *rport, Zone: ""})
	if err != nil {
		log.Errorf("can not create Conn %s", err.Error())
		return
	}
	defer Conn.Close()

	go ConnRecv(Conn)
	time.Now()
	for {
		select {
		case <-time.After(1 * time.Second):
			mu.Lock()
			if recvTime.After(sendTime) {
				diff := recvTime.Sub(sendTime)
				log.Infof("%v", diff)
			} else {
				log.Infof("lost")
			}
			sendTime = time.Now()
			mu.Unlock()
			Conn.Write([]byte("hello"))
		}
	}
}
