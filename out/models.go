package out

import (
	"github.com/miclip/nuget-resource"
)

type Request struct {
	Source nugetresource.Source `json:"source"`
	Params Params               `json:"params"`
}

type Params struct {
	PackagePath string `json:"package_path"`
}

type Response struct {
	Version  nugetresource.Version        `json:"version"`
	Metadata []nugetresource.MetadataPair `json:"metadata"`
}

