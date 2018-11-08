package check

import "github.com/miclip/nuget-resource"

type Request struct {
	Source  nugetresource.Source  `json:"source"`
	Version nugetresource.Version `json:"version"`
}

type Response []nugetresource.Version