package jisho

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	BaseUrl string
}

func New() *Client {
	return &Client{}
}

type SearchWordsResponse struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Data []struct {
		Slug     string `json:"slug"`
		Japanese []struct {
			Word    string `json:"word"`
			Reading string `json:"reading"`
		} `json:"japanese"`
		Senses []struct {
			EnglishDefinitions []string `json:"english_definitions"`
			PartsOfSpeech      []string `json:"parts_of_speech"`
			Tags               []string `json:"tags"`
		} `json:"senses"`
	} `json:"data"`
}

func (c *Client) SearchWords(keyword string) (*SearchWordsResponse, error) {
	u := &url.URL{
		Scheme:   "https",
		Host:     "jisho.org",
		Path:     fmt.Sprintf("api/v1/search/words"),
		RawQuery: fmt.Sprintf("keyword=%s", url.PathEscape(keyword)),
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Status %d: %s", resp.StatusCode, body)
	}
	swr := &SearchWordsResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(swr); err != nil {
		return nil, err
	}
	return swr, nil

}
