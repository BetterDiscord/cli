package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func DownloadFile(url string, filepath string) (response *http.Response, err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := out.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Setup the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "BetterDiscord/cli")
	req.Header.Add("Accept", "application/octet-stream")

	// Get the data
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("bad status code: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func DownloadJSON[T any](url string) (T, error) {
	var data T

	// Setup the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return data, err
	}
	req.Header.Add("User-Agent", "BetterDiscord/cli")

	// Get the data
	resp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("bad status: %s", resp.Status)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}
