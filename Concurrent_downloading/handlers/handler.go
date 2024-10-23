package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"sync"
)

type imageField struct {
	URL string `json:"url"` // Capitalized to be exported
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/fetch" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Only Get request are Acceptable", http.StatusBadRequest)
		return
	}

	apikey := r.Header.Get("api_key")
	if apikey == "" {
		http.Error(w, "API key is not provided", http.StatusBadRequest)
		return
	}

	apiList := []string{
		fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apikey),
		fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apikey),
		fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apikey),
		fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apikey),
		fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apikey),
		fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apikey),
		fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apikey),
		fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apikey),
		fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apikey),
	}
	results := make([]imageField, len(apiList))
	var wg sync.WaitGroup
	var mu sync.Mutex

	if err := os.MkdirAll("./images", os.ModePerm); err != nil {
		http.Error(w, "Failed to create images directory", http.StatusInternalServerError)
		return
	}
	for i, url := range apiList {
		wg.Add(1)
		go func(url string) {

			defer wg.Done()
			response, err := http.Get(url)
			if err != nil {
				http.Error(w, "Failed to fetch data from NASA API", http.StatusInternalServerError)
				return
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				http.Error(w, "Failed to read response from NASA API", http.StatusInternalServerError)
				return
			}

			var result imageField

			if err := json.Unmarshal(body, &result); err != nil {
				fmt.Println("Failed to unmarshal JSON:", err)
				return
			}

			fmt.Println(result.URL)
			imgResponce, err := http.Get(result.URL)
			if err != nil {
				http.Error(w, "Error getting Image", http.StatusInternalServerError)
				return
			}

			defer imgResponce.Body.Close()
			// Step 2: Read the response body into a byte slice
			imgBody, err := io.ReadAll(imgResponce.Body)
			if err != nil {
				fmt.Println("Failed to read image response body:", err)
				return
			}

			image, _, err := image.Decode(bytes.NewReader(imgBody))
			if err != nil {
				fmt.Println("Failed to decode image:", err)
				return
			}
			out, err := os.Create(fmt.Sprintf("./images/myimage_%d.jpg", i)) // Unique filename
			if err != nil {
				fmt.Println("Failed to create file:", err)
				return
			}
			defer out.Close()

			// Write the image to  the file as JPEG
			if err := jpeg.Encode(out, image, nil); err != nil {
				fmt.Println("Failed to save image:", err)
				return
			}

			mu.Lock()
			results[i] = result
			mu.Unlock()
		}(url)
	}
	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(results) // Marshal the results into JSON
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData) // Write the JSON response

}
