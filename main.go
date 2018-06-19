package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"gopkg.in/cheggaaa/pb.v1"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "   webpoll --attempts 10 --timeout 500 http://www.example.org\n")
	fmt.Fprintf(os.Stderr, "where\n")
	fmt.Fprintf(os.Stderr, "   timeout    the timeout before an attempt is considered failed\n")
	fmt.Fprintf(os.Stderr, "   attempts   the unsuccessful attempts before the command fails\n")
	os.Exit(1)
}

func main() {
	var timeout, attempts int
	flag.IntVar(&timeout, "timeout", 1000, "number of milliseconds between attempth [default: 1000]")
	flag.IntVar(&attempts, "attempts", 60, "number of attempts before failing [default: 60]")
	flag.Parse()

	if len(flag.Args()) == 0 {
		usage()
	}

	address := flag.Args()[0]

	fmt.Printf(" Connecting to %q (%d attempts, %dms each)\n", address, attempts, timeout)

	transport := http.Transport{
		Dial: func(network string, address string) (net.Conn, error) {
			return net.DialTimeout(network, address, time.Duration(timeout)*time.Millisecond)
		},
		TLSClientConfig: &tls.Config{
			// MaxVersion:         tls.VersionTLS11,
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{
		Transport: &transport,
	}

	bar := pb.StartNew(attempts)
	bar.ShowTimeLeft = false
	bar.ShowFinalTime = false
	bar.SetRefreshRate(20 * time.Millisecond)

	spinner := []byte{'|', '/', '-', '\\', '|', '/', '-', '\\'}
	start := time.Now()
	for attempt := 0; attempt < attempts; attempt++ {
		t0 := time.Now()
		response, err := client.Get(address)
		elapsed := time.Since(t0)
		bar.Increment()
		bar.BarStart = string(spinner[attempt%len(spinner)])
		if err == nil {
			response.Body.Close()
			if response.StatusCode >= 200 && response.StatusCode < 400 {
				bar.BarStart = "["
				bar.FinishPrint(fmt.Sprintf(" Connection succeeded after %d attempts (elapsed: %s)", attempt+1, time.Since(start)))
				os.Exit(0)
			}
		}
		if elapsed < (time.Duration(timeout) * time.Millisecond) {
			time.Sleep(time.Duration(timeout)*time.Millisecond - elapsed)
		}
	}
	bar.BarStart = "["
	bar.FinishPrint(fmt.Sprintf(" No connection to server after %d attempts (elapsed: %s)", attempts, time.Since(start)))
	os.Exit(1)
}
