package nuget

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/xmlpath.v1"
)

const (
	isPackableXpath  = "/Project/PropertyGroup/IsPackable"
	packageIDXpath   = "/Project/PropertyGroup/PackageId"
	authorsXpath     = "/Project/PropertyGroup/Authors"
	ownerXpath       = "/Project/PropertyGroup/Company"
	descriptionXpath = "/Project/PropertyGroup/Description"
	versionXpath     = "/Project/PropertyGroup/Version"
)

// NugetClientv3 ...
type NugetClientv3 interface {
	GetServiceIndex(ctx context.Context) (*ServiceIndex, error)
	SearchQueryService(ctx context.Context, searchQueryURL string, query string, preRelease bool) (*SearchResults, error)
	GetPackageVersions(ctx context.Context, name string, preRelease bool) ([]Version, error)
	CreateNuspec(packageID string, version string, author string, description string, owner string) Nuspec
	DownloadPackage(ctx context.Context, packageID string, version string, targetFolder string) error
	GetNugetApiEndPoint(ctx context.Context, resourceType string) (string, error)
	CreateNuspecFromProject(project string, version string) (Nuspec, error)
	AutoIncrementVersion(versionSpec string, version string) (string, error)
	PublishPackage(ctx context.Context, apikey string, packagePath string) error
}

type nugetclientv3 struct {
	FeedURL      string
	ServiceIndex ServiceIndex
}

func NewNugetClientv3(
	feedurl string,
) NugetClientv3 {

	return &nugetclientv3{
		FeedURL: feedurl,
	}
}

func (client *nugetclientv3) GetServiceIndex(ctx context.Context) (*ServiceIndex, error) {

	req, err := http.NewRequest(http.MethodGet, client.FeedURL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Add("accept", "application/json")
	var netClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error getting Service Index %d", res.StatusCode)
	}
	defer res.Body.Close()

	var r ServiceIndex
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (client *nugetclientv3) SearchQueryService(ctx context.Context, searchQueryURL string, query string, preRelease bool) (*SearchResults, error) {
	queryParams := fmt.Sprintf("?q=%s&prerelease=%t", query, preRelease)
	req, err := http.NewRequest(http.MethodGet, searchQueryURL+queryParams, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Add("accept", "application/json")
	var netClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error querying Service %d", res.StatusCode)
	}
	defer res.Body.Close()

	var r SearchResults
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (client *nugetclientv3) PublishPackage(ctx context.Context, apikey string, packagePath string) error {
	uploadURL, err := client.GetNugetApiEndPoint(ctx, "PackagePublish/2.0.0")
	if err != nil {
		log.Fatal("error getting download url", err)
	}
	file, err := os.Open(packagePath)
	if err != nil {
		return fmt.Errorf("error reading package %s with %v", packagePath, err)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("package", file.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, uploadURL, body)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Add("X-NuGet-ApiKey", apikey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	var netClient = &http.Client{
		Timeout: 300 * time.Second,
	}
	res, err := netClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 && res.StatusCode != 201 && res.StatusCode != 202 {
		return fmt.Errorf("error uploading package %d", res.StatusCode)
	}
	defer res.Body.Close()

	return nil
}

func (client *nugetclientv3) DownloadPackage(ctx context.Context, packageID string, version string, targetFolder string) error {

	downloadURL, err := client.GetNugetApiEndPoint(ctx, "PackageBaseAddress/3.0.0")
	if err != nil {
		log.Fatal("error getting download url", err)
	}

	targetFolder = targetFolder + "/packages"

	err = os.MkdirAll(targetFolder, 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(targetFolder + "/" + packageID + "." + version + ".nupkg")
	if err != nil {
		return err
	}
	defer out.Close()

	queryParams := fmt.Sprintf("%s/%s/%s", packageID, version, packageID+"."+version+".nupkg")
	req, err := http.NewRequest(http.MethodGet, downloadURL+queryParams, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	var netClient = &http.Client{
		Timeout: 300 * time.Second,
	}
	res, err := netClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("error downloading package %d", res.StatusCode)
	}
	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func (client *nugetclientv3) CreateNuspec(packageID string, version string, author string, description string, owner string) Nuspec {
	return Nuspec{
		Xmlns:                    "http://schemas.microsoft.com/packaging/2013/05/nuspec.xsd",
		ID:                       packageID,
		Version:                  version,
		Authors:                  author,
		Owners:                   owner,
		RequireLicenseAcceptance: false,
		Description:              description,
	}
}

func (client *nugetclientv3) CreateNuspecFromProject(project string, version string) (Nuspec, error) {
	isPackablePath := xmlpath.MustCompile(isPackableXpath)
	packageIDPath := xmlpath.MustCompile(packageIDXpath)
	authorPath := xmlpath.MustCompile(authorsXpath)
	ownerPath := xmlpath.MustCompile(ownerXpath)
	versionPath := xmlpath.MustCompile(versionXpath)
	descriptionPath := xmlpath.MustCompile(descriptionXpath)

	file, err := os.Open(project)
	if err != nil {
		return Nuspec{}, err
	}
	root, err := xmlpath.Parse(file)
	if err != nil {
		return Nuspec{}, err
	}
	if isPackable, ok := isPackablePath.String(root); ok {
		if isPackable == "true" {
			var id, authors, owners, description, defaultVersion string
			if value, ok := packageIDPath.String(root); ok {
				id = value

				if value, ok := authorPath.String(root); ok {
					authors = value
				}
				if value, ok := ownerPath.String(root); ok {
					owners = value
				}
				if value, ok := descriptionPath.String(root); ok {
					description = value
				}
				if value, ok := versionPath.String(root); ok {
					defaultVersion = value
				}
				if version == "" {
					version = defaultVersion
				}
				return client.CreateNuspec(id, version, authors, description, owners), nil
			}
		}
	}
	return Nuspec{}, fmt.Errorf("the project file could not be parsed %s", project)
}

func (client *nugetclientv3) GetNugetApiEndPoint(ctx context.Context, resourceType string) (string, error) {
	serviceIndex, err := client.GetServiceIndex(ctx)
	if err != nil {
		return "", err
	}

	for _, resource := range serviceIndex.Resources {
		if resource.Type == resourceType {
			return resource.ID, nil
		}
	}

	return "", fmt.Errorf("Could not find %s Endpoint", resourceType)

}

func (client *nugetclientv3) GetPackageVersions(ctx context.Context, name string, preRelease bool) ([]Version, error) {

	versions := []Version{}

	searchQueryService, err := client.GetNugetApiEndPoint(ctx, "SearchQueryService")
	if err != nil {
		log.Fatal("error getting package version", err)
	}

	searchResults, err := client.SearchQueryService(ctx, searchQueryService, name, preRelease)
	if err != nil {
		return versions, err
	}

	if searchResults == nil {
		return versions, fmt.Errorf("Package not found name: %s prerelease: %t", name, preRelease)
	}

	for _, result := range searchResults.Data {
		if result.ID == name {
			return result.Versions, nil
		}
	}

	return nil, nil

}

func (client *nugetclientv3) AutoIncrementVersion(versionSpec string, version string) (string, error) {
	if versionSpec == "" {
		return version, nil
	}
	latestVersion := strings.Split(version, ".")
	specVersion := strings.Split(versionSpec, ".")
	if len(latestVersion) != len(specVersion) {
		return "", fmt.Errorf("Version semantics don't match version spec: %s version: %s ", versionSpec, version)
	}
	for index := 0; index < len(specVersion); index++ {
		if specVersion[index] == "*" {
			i, _ := strconv.Atoi(latestVersion[index])
			latestVersion[index] = strconv.Itoa(i + 1)
		}
	}
	return strings.Join(latestVersion, "."), nil
}
