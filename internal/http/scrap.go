package http

import (
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

func ScrapURL(url string) (page string, err error) {
	var (
		res     *http.Response
		content []byte
	)
	res, err = http.Get(url)
	if err != nil {
		log.Err(err).Msgf("Failed to fetch URL: %s", url)
		return
	}
	content, err = io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		log.Err(err).Msg("Unable to extract body content")
		return
	}

	page = string(content)
	return
}
