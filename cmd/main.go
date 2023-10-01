package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/yomorun/yomo"
	"github.com/yomorun/yomo/serverless"
)

const (
	INPUT_TAG  = 0x33
	OUTPUT_TAG = 0x34

	PREFIX = "\u001B[32m[cli]\u001B[0m"
)

var ch chan *ImageResult

type ImageResult struct {
	Score float32 `json:"score"`
	Class int32   `json:"class"`
}

func Handler(ctx serverless.Context) {
	data := ctx.Data()
	var result ImageResult
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalln(err)
	}

	ch <- &result
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: cli image")
		return
	}
	inputImage := os.Args[1]

	ch = make(chan *ImageResult)

	addr := "localhost:9000"
	if v := os.Getenv("YOMO_ADDR"); v != "" {
		addr = v
	}
	sfn := yomo.NewStreamFunction("sink", addr)
	sfn.SetHandler(Handler)
	sfn.SetObserveDataTags(OUTPUT_TAG)
	err := sfn.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer sfn.Close()

	source := yomo.NewSource("source", addr)
	err = source.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer source.Close()

	image, err := os.ReadFile(inputImage)
	if err != nil {
		log.Fatalln(err)
	}
	source.Write(INPUT_TAG, image)
	fmt.Printf("%s send image: %s\n", PREFIX, inputImage)

	result := <-ch
	fmt.Printf("%s inference result: score=%f, class=%d\n", PREFIX, result.Score, result.Class)
}
