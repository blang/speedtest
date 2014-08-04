package speedtest

import (
	"crypto/rand"
	"fmt"
	"github.com/dustin/go-humanize"
	"log"
	"math"
	"net"
	"time"
)

type BytesPerTime struct {
	Bytes    uint64
	Duration time.Duration
}

func SpeedMeter(input chan BytesPerTime, bytesPerSec chan BytesPerTime) {
	go func() {
		bpt := BytesPerTime{}
		for {
			select {
			case newBpt := <-input:
				bpt.Bytes += newBpt.Bytes
				bpt.Duration += newBpt.Duration
			case bytesPerSec <- bpt:
				bpt = BytesPerTime{}
			}
		}
	}()
}

func SpeedReporter(input chan BytesPerTime, interval time.Duration) {
	go func() {
		for {
			select {
			case <-time.After(interval):
				bpt, ok := <-input
				if !ok {
					log.Println("Reporting stopped")
					return
				}
				if bpt.Duration.Seconds() != 0 {
					log.Printf("Throughput: %s/s", humanize.IBytes(uint64(math.Ceil(float64(bpt.Bytes)/bpt.Duration.Seconds()))))
				} else {
					log.Println("No throughput")
				}
			}
		}
	}()
}

func SendData(conn net.Conn, buffersize int, reportCh chan BytesPerTime) error {
	buffer := make([]byte, buffersize)
	read, err := rand.Read(buffer)
	if err != nil {
		return fmt.Errorf("Error while initialising buffer: %s", err)
	}
	if read != buffersize {
		return fmt.Errorf("Could not init buffer: %d read", read)
	}

	for {
		startTime := time.Now()

		w, err := conn.Write(buffer)
		if err != nil {
			return fmt.Errorf("Error while writing: %s", err)
		}

		reportCh <- BytesPerTime{
			Bytes:    uint64(w),
			Duration: time.Since(startTime),
		}
	}
}

func ReceiveData(conn net.Conn, buffersize int, reportCh chan BytesPerTime) error {
	b := make([]byte, buffersize)
	defer conn.Close()
	for {
		startTime := time.Now()

		w, err := conn.Read(b)

		if err != nil {
			return fmt.Errorf("Read: %d, Error: %s\n", w, err)
		}

		reportCh <- BytesPerTime{
			Bytes:    uint64(w),
			Duration: time.Since(startTime),
		}
	}
}
