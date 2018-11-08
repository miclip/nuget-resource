package in

import (
	"context"
	"os"
	"github.com/miclip/nuget-resource/nuget"
	"github.com/miclip/nuget-resource"
)

//Execute - provides in capability
func Execute(request Request, targetDir string) (Response, []byte, error) {
	out := []byte{}

	if request.Version.PackageID == "" {
		return Response{}, out, nil
	}

	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return Response{}, out, err
	}
	
	nugetclient := nuget.NewNugetClientv3(request.Source.NugetSource)
	err = nugetclient.DownloadPackage(context.Background(), request.Version.PackageID, request.Version.Version, targetDir)
	if err != nil {
		return Response{}, out, err
	}
	nugetresource.Sayf("downloaded package %s %s \n",request.Version.PackageID, request.Version.Version)

	

	return Response{}, out, nil
}