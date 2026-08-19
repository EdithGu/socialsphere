package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"socialai/service"
)

type generateImageRequestBody struct {
	Prompt string `json:"prompt"`
}

type generateImageResponseBody struct {
	Url string `json:"url"`
}

func generateImageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one image generation request")
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var reqBody generateImageRequestBody
	if err := decoder.Decode(&reqBody); err != nil {
		http.Error(w, "Cannot decode request body", http.StatusBadRequest)
		fmt.Printf("Cannot decode request body %v\n", err)
		return
	}

	if reqBody.Prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	imageUrl, err := service.GenerateImage(reqBody.Prompt)
	if err != nil {
		http.Error(w, "Failed to generate image", http.StatusInternalServerError)
		fmt.Printf("Failed to generate image %v\n", err)
		return
	}

	js, err := json.Marshal(generateImageResponseBody{Url: imageUrl})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
