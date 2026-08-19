package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"socialai/constants"
)

const openAIImagesURL = "https://api.openai.com/v1/images/generations"

type imageGenerationRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type imageGenerationResponse struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

// GenerateImage calls OpenAI's DALL-E 3 image generation API from the
// backend, so the OpenAI API key stays server-side (read from the
// OPENAI_API_KEY environment variable) instead of shipping in the
// frontend JS bundle.
func GenerateImage(prompt string) (string, error) {
	if constants.OPENAI_API_KEY == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is not configured")
	}

	reqBody, err := json.Marshal(imageGenerationRequest{
		Model:  "dall-e-3",
		Prompt: prompt,
		N:      1,
		Size:   "1024x1024",
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, openAIImagesURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+constants.OPENAI_API_KEY)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI image generation failed (%d): %s", resp.StatusCode, string(body))
	}

	var parsed imageGenerationResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Data) == 0 || parsed.Data[0].URL == "" {
		return "", fmt.Errorf("OpenAI response did not include an image URL")
	}

	return parsed.Data[0].URL, nil
}
