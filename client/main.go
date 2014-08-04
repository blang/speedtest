package main

import (
	"flag"
	"github.com/blang/speedtest"
	"log"
	"net"
	"time"
)

func main() {
	connect := flag.String("connect", "", "connect to address")
	buffersize := flag.Int("buffer", 4096, "Buffer size")
	reportinterval := flag.Duration("report", 5*time.Second, "Report interval")
	send := flag.Bool("send", true, "True for send, false for receive")
	flag.Parse()

	conn, err := net.Dial("tcp", *connect)
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	reportCh := make(chan speedtest.BytesPerTime)
	statsCh := make(chan speedtest.BytesPerTime)
	speedtest.SpeedMeter(reportCh, statsCh) // Speedmeter on all connections
	speedtest.SpeedReporter(statsCh, *reportinterval)

	if *send { // Client send mode
		log.Println("Enter Send mode")
		err := speedtest.SendData(conn, *buffersize, reportCh)
		if err != nil {
			log.Printf("Error: %s", err)
			return
		}
	} else { // Receive mode
		log.Println("Enter Receive mode")
		err := speedtest.ReceiveData(conn, *buffersize, reportCh)
		if err != nil {
			log.Printf("Error: %s", err)
			return
		}
	}

}
