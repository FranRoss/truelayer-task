package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// models, in this file for simplicity, could go in a separate file
type TranslationType string

const (
	Yoda        TranslationType = "yoda.json"
	Shakespeare TranslationType = "shakespeare.json"
)

type TranslationRequest struct {
	Text string `json:"text"`
}

type TranslationResponse struct {
	Success struct {
		Total int `json:"total"`
	} `json:"success"`
	Contents struct {
		Translated  string `json:"translated"`
		Text        string `json:"text"`
		Translation string `json:"translation"`
	} `json:"contents"`
	Error *APIErrorDetails `json:"error,omitempty"`
}

type APIErrorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// client
type TranslationClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewTranslationClient() *TranslationClient {
	return &TranslationClient{
		baseURL: "https://api.funtranslations.mercxry.me/v1",
		httpClient: &http.Client{
			Timeout: CLIENTS_TIMEOUT * time.Second,
		},
	}
}

func (c *TranslationClient) Translate(text string, translationType TranslationType) (string, error) {
	url := fmt.Sprintf("%s/translate/%s", c.baseURL, translationType)
	payload := TranslationRequest{Text: text}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("failed making request to FunTranslations: %w", err)
	}
	defer resp.Body.Close()

	// this handles the 429 for rate limiting
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http error status %d", resp.StatusCode)
	}

	res := &TranslationResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed decoding JSON response: %w", err)
	}

	if res.Error != nil {
		return "", fmt.Errorf("%d %s", res.Error.Code, res.Error.Message)
	}

	return res.Contents.Translated, nil
}
