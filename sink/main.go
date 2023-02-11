package main

import (
	"encoding/json"
	"fmt"

	"github.com/yomorun/yomo/core/frame"
)

type ImageResult struct {
	Score float32 `json:"score"`
	Class int32   `json:"class"`
}

func DataTags() []frame.Tag {
	return []frame.Tag{0x34}
}

func Handler(data []byte) (frame.Tag, []byte) {
	var result ImageResult
	err := json.Unmarshal(data, &result)
	if err != nil {
		fmt.Errorf("json.Unmarshal error: %v\n", err)
		return 0x0, nil
	}

	fmt.Printf("score: %f, class: %d\n", result.Score, result.Class)
	return 0x0, nil
}
