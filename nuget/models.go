package nuget

import (
	"encoding/xml"
)

// ServiceIndex response type
type ServiceIndex struct {
	Version   string     `json:"version"`
	Resources []Resource `json:"resources"`
}

// Resource nuget Resource type
type Resource struct {
	ID      string `json:"@id"`
	Type    string `json:"@type"`
	Comment string `json:"comment"`
}

// SearchResults from the nuget api
type SearchResults struct {
	TotalHits  int            `json:"totalHits"`
	Index      string         `json:"index"`
	LastReopen string         `json:"lastReopen"`
	Data       []SearchResult `json:"data"`
}

// SearchResult from the nuget api
type SearchResult struct {
	ID          string    `json:"id"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Versions    []Version `json:"versions"`
}

type PackageVersion struct {
	ID          string `json:"id"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type Version struct {
	ID        string `json:"id"`
	Version   string `json:"version"`
	Downloads int    `json:"downloads"`
}

type Nuspec struct {
	XMLName                  xml.Name `xml:"package"`
	Xmlns                    string   `xml:"xmlns,attr"`
	ID                       string   `xml:"metadata>id"`
	Version                  string   `xml:"metadata>version"`
	Authors                  string   `xml:"metadata>authors"`
	Owners                   string   `xml:"metadata>owners"`
	RequireLicenseAcceptance bool     `xml:"metadata>requireLicenseAcceptance"`
	Description              string   `xml:"metadata>description"`
}
