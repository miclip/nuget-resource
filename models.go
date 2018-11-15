package nugetresource

type Source struct {
	NugetSource  string `json:"nuget_source"`
	NugetAPIKey  string `json:"nuget_apikey"`
	NugetTimeout int    `json:"nuget_timeout"`
	PackageID    string `json:"package_Id"`
	PreRelease   bool   `json:"prerelease"`
}

func (source Source) IsValid() (bool, string) {

	return true, ""
}

type Version struct {
	PackageID string `json:"package_id"`
	Version   string `json:"version"`
}

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
