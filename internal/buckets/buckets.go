package buckets

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

type BunnyNetFile struct {
	Guid            string `json:"Guid"`
	StorageZoneName string `json:"StorageZoneName"`
	Path            string `json:"Path"`
	ObjectName      string `json:"ObjectName"`
	Length          int    `json:"Length"`
	LastChanged     string `json:"LastChanged"`
	ServerId        int    `json:"ServerId"`
	ArrayNumber     int    `json:"ArrayNumber"`
	IsDirectory     bool   `json:"IsDirectory"`
	UserId          string `json:"UserId"`
	ContentType     string `json:"ContentType"`
	DateCreated     string `json:"DateCreated"`
	StorageZoneId   int    `json:"StorageZoneId"`
	Checksum        string `json:"Checksum"`
	ReplicatedZones string `json:"ReplicatedZones"`
}

const (
	storageZone = "zzep"
)

func GetFile(fileName string) ([]byte, error) {
	apiKey := os.Getenv("FILE_STORAGE_PASSWORD")

	// Construct the URL
	url := fmt.Sprintf("https://storage.bunnycdn.com/%s/%s", storageZone, fileName)
	if os.Getenv("NODE_ENV") == "development" {
		url = fmt.Sprintf("https://storage.bunnycdn.com/%s/development/%s", storageZone, fileName)
	}

	// Create the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("AccessKey", apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	return body, nil
}

func UploadFile(filePath string, fileName string) (string, error) {
	apiKey := os.Getenv("FILE_STORAGE_PASSWORD")

	// Read the file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Construct the URL
	url := fmt.Sprintf("https://storage.bunnycdn.com/%s/%s", storageZone, fileName)

	// Create the PUT request
	req, err := http.NewRequest("PUT", url, bytes.NewReader(fileData))
	if err != nil {
		return "", err
	}

	// Add headers
	req.Header.Set("AccessKey", apiKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode == 201 {
		return fmt.Sprintf("https://zzep.b-cdn.net/%s", fileName), err
	} else {
		return "", nil
	}
}

func DeleteFile(fileName string) (bool, error) {
	apiKey := os.Getenv("FILE_STORAGE_PASSWORD")

	// Construct the URL
	url := fmt.Sprintf("https://storage.bunnycdn.com/%s/%s", storageZone, fileName)
	if os.Getenv("NODE_ENV") == "development" {
		url = fmt.Sprintf("https://storage.bunnycdn.com/%s/development/%s", storageZone, fileName)
	}

	// Create the DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return false, err
	}

	// Add headers
	req.Header.Set("AccessKey", apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode == 200 {
		return true, err
	} else {
		return false, err
	}
}

func DeleteFolder(folderName string) (bool, error) {
	apiKey := os.Getenv("FILE_STORAGE_PASSWORD")

	// Construct the URL
	url := fmt.Sprintf("https://storage.bunnycdn.com/%s/%s/", storageZone, folderName)
	if os.Getenv("NODE_ENV") == "development" {
		url = fmt.Sprintf("https://storage.bunnycdn.com/%s/development/%s/", storageZone, folderName)
	}

	// Create the DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return false, err
	}

	// Add headers
	req.Header.Set("AccessKey", apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode == 200 {
		return true, err
	} else {
		return false, err
	}
}

func ListFiles(folderName string) (string, error) {
	apiKey := os.Getenv("FILE_STORAGE_PASSWORD")

	// Construct the URL
	url := fmt.Sprintf("https://storage.bunnycdn.com/%s/%s", storageZone, folderName)
	if os.Getenv("NODE_ENV") == "development" {
		url = fmt.Sprintf("https://storage.bunnycdn.com/%s/development/%s", storageZone, folderName)
	}

	// Create the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Add headers
	req.Header.Set("AccessKey", apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	return string(body), nil
}
