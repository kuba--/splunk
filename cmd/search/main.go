package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kuba--/splunk"
)

var info = flag.Bool("info", false, "The server's info.")
var from = flag.String("from", "-5min", "The earliest time for the time range of your search.")

type jsonWriter struct {
	buf bytes.Buffer
}

func (w *jsonWriter) Write(data []byte) (int, error) {
	json.Indent(&w.buf, data, "", "\t")
	w.buf.WriteTo(os.Stdout)

	return w.buf.Len(), nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{}, 1)
	defer func() {
		done <- struct{}{}
		close(done)
	}()

	sig := make(chan os.Signal, 1)
	// Handle SIGINT and SIGTERM.
	signal.Notify(sig, os.Interrupt, os.Kill)
	go func() {
		select {
		case <-sig:
			time.AfterFunc(2*time.Second, func() { os.Exit(0) })
			cancel()

		case <-done:
		}
	}()

	client := splunk.NewClient(os.Getenv("SPLUNK_USERNAME"), os.Getenv("SPLUNK_PASSWORD"), os.Getenv("SPLUNK_URL"))
	if err := client.Login(ctx); err != nil {
		log.Fatalln(err)
	}

	w := &jsonWriter{}
	if *info {
		if err := client.Info(ctx, w); err != nil {
			log.Fatalln(err)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if q := scanner.Text(); q != "" {
				if err := client.Search(ctx, q, *from, w); err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
}
