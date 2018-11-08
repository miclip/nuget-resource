package out

import (
	"time"
	"os/exec"

	"github.com/miclip/nuget-resource"
)

var ExecCommand = exec.Command

//Execute - provides out capability
func Execute(request Request, sourceDir string) (Response, []byte, error) {
	out := []byte{}
	response := Response{
		Version: nugetresource.VersionTime{
			Timestamp: time.Now(),
		},
		Metadata: []nugetresource.MetadataPair{
			{
				Name:  "project",
				Value: request.Params.Project,
			},
			{
				Name:  "framework",
				Value: request.Source.Framework,
			},
			{
				Name:  "runtime",
				Value: request.Source.Runtime,
			},
		},
	}



	return response, out, nil
}


