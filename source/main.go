package main

import (
	"log"
	"os"
	"time"

	"github.com/yomorun/yomo"
)

func main() {
	addr := "localhost:9000"
	if v := os.Getenv("YOMO_ADDR"); v != "" {
		addr = v
	}

	source := yomo.NewSource(
		"source",
		yomo.WithZipperAddr(addr),
	)
	source.SetDataTag(0x33)
	if err := source.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer source.Close()

	image, err := os.ReadFile("./sample.png")
	if err != nil {
		log.Fatalln(err)
	}
	source.Write(image)
	time.Sleep(3 * time.Second) // wait for sending data
}
