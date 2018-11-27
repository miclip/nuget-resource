package in

import (
	"context"
	"os"

	"github.com/miclip/nuget-resource"
	"github.com/miclip/nuget-resource/nuget"
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
	file, err := nugetclient.DownloadPackage(context.Background(), request.Version.PackageID, request.Version.Version, targetDir)
	if err != nil {
		return Response{}, out, err
	}
	nugetresource.Sayf("downloaded package %s %s \n", request.Version.PackageID, request.Version.Version)

	err = nugetresource.UnarchiveZip(file, targetDir)
	if err != nil {
		nugetresource.Fatal("error extracting package", err)
	}
	nugetresource.Sayf("extracted archive %s to %s", file, targetDir)

	nugetresource.ChmodAllFiles(targetDir, 0600)

	response := Response{
		Version : request.Version,
		Metadata : []nugetresource.MetadataPair{},		
	}	

	return response, out, nil
}
