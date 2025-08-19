package util

import "net/http"

func MeilisearchHeaders(header http.Header) (host, apiKey, hash string) {
	host = header.Get("X-Meili-Instance")
	apiKey = header.Get("X-Meili-APIKey")

	hash = Hash(host + ":" + apiKey)
	return
}
