package main

import (
    "log"
    "testing"
)

func TestSubtitle(t *testing.T) {
    err := SubtitleProcess(`https://i0.hdslb.com/bfs/subtitle/3ec566113d8d6e525919ff4efadc2fa7895f26a1.json`, `test`)
    if err != nil {
        log.Println(err)
    }
}
