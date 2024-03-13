package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

var (
	cacheMutex sync.Mutex
	cache      = make(map[int]ChapterCache)
)

type ChapterCache struct {
	Meaning string
	Summary string
}

func GetChapterSummary(chapterNumber int) (string, string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if cached, ok := cache[chapterNumber]; ok {
		return cached.Meaning, cached.Summary
	}

	resp, err := http.Get(fmt.Sprintf("https://bhagavadgitaapi.in/chapter/%d", chapterNumber))
	if err != nil {
		log.Println(err)
		return "", ""
	}
	defer resp.Body.Close()

	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal(jsonResponse, &data); err != nil {
		log.Printf("Error unmarshaling JSON, %s", err)
		return "", ""
	}

	meaning := data["meaning"].(map[string]interface{})
	summary := data["summary"].(map[string]interface{})

	meaningEn := meaning["en"].(string)
	summaryEn := summary["en"].(string)

	// Cache the retrieved data
	cache[chapterNumber] = ChapterCache{Meaning: meaningEn, Summary: summaryEn}

	return meaningEn, summaryEn
}

func main() {
	chapter := rand.Intn(18)
	meaning, summary := GetChapterSummary(chapter)
	fmt.Printf("Chapter %d:\nMeaning: %s\nSummary: %s\n", chapter, meaning, summary)
}
