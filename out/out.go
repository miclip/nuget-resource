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

	// sleep to allow time for nuget caches to clear 
	if request.Params.NugetCacheDelay > 0 {
		nugetresource.Sayf("waiting for cache refresh, %v seconds", request.Params.NugetCacheDelay)
		time.Sleep(time.Duration(request.Params.NugetCacheDelay) * time.Millisecond)
	}	

	packageVersions, err := nugetclient.GetPackageVersions(context.Background(), request.Source.PackageID, request.Source.PreRelease)
	if err != nil {
		nugetresource.Fatal("error querying for latest version from nuget", err)
	}

	var version string
	if len(packageVersions)==0 {
		version = "?"
	} else {
		version = packageVersions[len(packageVersions)-1].Version
	}

	response := Response{
		Version: nugetresource.Version{
			PackageID: request.Source.PackageID,
			Version: version,
		},
		Metadata: []nugetresource.MetadataPair{},
	}
	return response, out, nil
}


