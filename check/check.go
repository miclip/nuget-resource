package check

import (
	"context"

	"github.com/miclip/nuget-resource"
	"github.com/miclip/nuget-resource/nuget"
)

//Execute - provides check capability
func Execute(request Request) (Response, error) {

	nugetclient := nuget.NewNugetClientv3(request.Source.NugetSource)
	packageVersions, err := nugetclient.GetPackageVersions(context.Background(), request.Source.PackageID, request.Source.PreRelease)
	if err != nil {
		nugetresource.Fatal("error querying for latest version from nuget.", err)
	}

	if packageVersions == nil {
		nugetresource.Sayf("package %s not found at %s ", request.Source.PackageID, request.Source.NugetSource)
		return Response{}, nil
	}

	response := []nugetresource.Version{}
	for _, pv := range packageVersions {
		response = append(response, nugetresource.Version{
			PackageID: request.Source.PackageID,
			Version:   pv.Version,
		})
	}
	
	return response, nil
}
