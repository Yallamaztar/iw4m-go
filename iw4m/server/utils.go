package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func readBody(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	return body, nil
}

func getDocFromRes(res *http.Response) (*goquery.Document, error) {
	body, err := readBody(res)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	return doc, nil
}

func (s *Server) getDoc(endpoint string) (*goquery.Document, error) {
	res, err := s.iw4m.DoRequest(endpoint)
	if err != nil {
		return nil, err
	}
	return getDocFromRes(res)
}
