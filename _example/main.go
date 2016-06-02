package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/tj/go-kinesis"
)

func main() {
	log.SetHandler(text.New(os.Stderr))
	log.SetLevel(log.DebugLevel)

	producer := kinesis.New(kinesis.Config{
		StreamName:  "logs",
		BacklogSize: 2000,
	})

	producer.Start()

	e := log.Entry{
		Message: "upload",
		Level:   log.InfoLevel,
		Fields: log.Fields{
			"app":  "site",
			"name": "sloth.png",
			"type": "image/png",
			"user": "Tobi",
		},
	}

	b, err := json.Marshal(e)
	if err != nil {
		log.WithError(err).Fatal("error marshalling")
	}

	go func() {
		for i := 0; i < 5000; i++ {
			err := producer.Put(b, "site")
			if err != nil {
				log.WithError(err).Fatal("error producing")
			}
		}
	}()

	time.Sleep(3 * time.Second)
	producer.Stop()
}
