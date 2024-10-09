package http

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"net/http"
)

func ScrapURL(url string) (page string, err error) {
	var (
		res *http.Response
		doc *goquery.Document
	)
	res, err = http.Get(url)
	if err != nil {
		log.Err(err).Msgf("Failed to fetch URL: %s", url)
		return
	}
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		log.Error().Err(err).Msg("Failed to fetch URL")
		return
	}

	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Err(err).Msg("Unable to create document from response")
		return
	}

	page = doc.Text()
	return
}
