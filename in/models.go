package in

import "github.com/miclip/nuget-resource"

type Request struct {
	Source  nugetresource.Source  `json:"source"`
	Version nugetresource.Version `json:"version"`
}

type Response struct {
	Version  nugetresource.Version        `json:"version"`
	Metadata []nugetresource.MetadataPair `json:"metadata"`
}