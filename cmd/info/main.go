package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kuba--/splunk"
)

type jsonWriter struct {
	buf bytes.Buffer
}

func (w *jsonWriter) Write(data []byte) (int, error) {
	json.Indent(&w.buf, data, "", "\t")
	w.buf.WriteTo(os.Stdout)

	return w.buf.Len(), nil
}

func main() {
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
	if err := client.Info(ctx, w); err != nil {
		log.Fatalln(err)
	}
}
