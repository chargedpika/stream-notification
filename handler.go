package main

import (
	"encoding/json"
	"net/http"
)

const CACHE_DURATION = 120

// Requests Twitch streams, answering with a filtered set.
func streamHandler(w http.ResponseWriter, r *http.Request) {
	// Always return JSON.
	w.Header().Set("Content-Type", "application/json")

	// Try getting a cached response.
	if cached := GetCache("streams"); cached != nil {
		w.Write(cached)
		return
	}

	sr, err := GetTwitchStreams(channels)
	if err != nil {
		w.Write([]byte("[]"))
		return
	}

	var resitems []ResponseItem
	for _, stream := range sr.Streams {
		if !gameFilter.MatchString(stream.Game) {
			continue
		}
		resitems = append(resitems, ResponseItem{Name: stream.Channel.DisplayName, URL: stream.Channel.URL})
	}

	b, err := json.Marshal(resitems)
	if err != nil || len(resitems) == 0 {
		w.Write([]byte("[]"))
		return
	}
	w.Write(b)
	PutCache("streams", CACHE_DURATION, b)
}
