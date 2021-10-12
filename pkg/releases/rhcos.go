package releases

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Generated with https://mholt.github.io/json-to-go/
type CincinnatiResponse struct {
	OCPReleases []OCPReleaseResponse `json:"nodes"`
}
type OCPMetaResponse struct {
	URL string `json:"url"`
}
type OCPReleaseResponse struct {
	Version  string          `json:"version"`
	Payload  string          `json:"payload"`
	Metadata OCPMetaResponse `json:"metadata"`
}

type OpenShiftVersion struct {
	Version string `json:"version" gorm:"unique_index"`
	Errata  string `json:"errata"`
	Image   string `json:"image"`
}

// Interface for getting Cincinnati data from URL or (test) file
type ReleaseDownloader func(path string) ([]byte, error)

// Get Cincinnati data from URL
func GetPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Got unexpected status code from proddefs endpoint: %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read resp.Body from proddefs endpoint")
	}
	return body, nil
}

func getEnvOrDefault(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		value = def
	}
	return value
}

func ParseCincinnati(path string, downloader ReleaseDownloader) ([]OpenShiftVersion, error) {
	byteValue, err := downloader(path)
	if err != nil {
		return nil, err
	}
	return unmarshallCincinnati(byteValue)
}

func unmarshallCincinnati(data []byte) ([]OpenShiftVersion, error) {
	var ocpVersions []OpenShiftVersion
	var response CincinnatiResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		fmt.Printf("Failed to marshal cincinnati response as JSON")
		return nil, err
	}
	for _, release := range response.OCPReleases {
		ocpVersion := OpenShiftVersion{
			Version: release.Version,
			Image:   release.Payload,
			Errata:  release.Metadata.URL,
		}
		ocpVersions = append(ocpVersions, ocpVersion)
	}
	return ocpVersions, nil
}
