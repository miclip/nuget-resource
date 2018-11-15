package out

import (
	"path"
	"time"
	"context"
	"github.com/miclip/nuget-resource/nuget"

	"github.com/miclip/nuget-resource"
)

//Execute - provides out capability
func Execute(request Request, sourceDir string) (Response, []byte, error) {
	out := []byte{}

	nugetclient := nuget.NewNugetClientv3(request.Source.NugetSource)
	err:=nugetclient.PublishPackage(context.Background(), request.Source.NugetAPIKey, path.Join(sourceDir,request.Params.PackagePath))
	if err != nil {
		nugetresource.Fatal("error publishing package to feed", err)
	}

	packageVersions, err := nugetclient.GetPackageVersions(context.Background(), request.Source.PackageID, request.Source.PreRelease)
	if err != nil {
		nugetresource.Fatal("error querying for latest version from nuget", err)
	}

	response := Response{
		Version: nugetresource.Version{
			PackageID: request.Source.PackageID,
			Version: packageVersions[len(packageVersions)-1].Version,
		},
		Metadata: []nugetresource.MetadataPair{
			nugetresource.MetadataPair{
			Name: request.Source.PackageID,
			Value: time.Now().String(),
			},
		},
	}
	return response, out, nil
}


