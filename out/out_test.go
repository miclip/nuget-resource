package out_test

import (
	"net/http"

	//"github.com/miclip/nuget-resource"
	"github.com/miclip/nuget-resource/nuget"
	//"github.com/miclip/nuget-resource/out"
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("out", func() {
	var server *ghttp.Server
	var returnedServiceIndex nuget.ServiceIndex
	var returnedSearchResults nuget.SearchResults
	var statusCode int
	BeforeEach(func() {
		server = ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/somefeed/api/v3/index.json"),
				ghttp.RespondWithJSONEncodedPtr(&statusCode, &returnedServiceIndex),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/somefeed/api/v2/package"),
				ghttp.VerifyHeader(http.Header{"X-NuGet-ApiKey": []string{"an_api_key"}}),
				ghttp.RespondWith(http.StatusOK,nil),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/somefeed/api/v3/query"),
				ghttp.RespondWithJSONEncodedPtr(&statusCode, &returnedSearchResults),
			),
		)
	})

	AfterEach(func() {
		server.Close()
	})

	Context("When the client returns a PackageVersion", func() {
		BeforeEach(func() {
			returnedServiceIndex = nuget.ServiceIndex{
				Version: "3.0.0",
				Resources: []nuget.Resource{
					nuget.Resource{
						ID:      server.URL() + "/somefeed/api/v3/query",
						Type:    "SearchQueryService",
						Comment: "Query endpoint of NuGet Search service.",
					},
					nuget.Resource{
						ID:      server.URL() + "/somefeed/api/v2/package",
						Type:    "PackagePublish/2.0.0",
						Comment: "Legacy gallery publish endpoint using the V2 protocol.",
					},
				},
			}
			returnedSearchResults = nuget.SearchResults{
				TotalHits:  1000,
				Index:      "index",
				LastReopen: "2018-10-22T22:45:00.425508Z",
				Data: []nuget.SearchResult{
					nuget.SearchResult{
						ID:          "Some.Package.Name",
						Version:     "2.0.2",
						Description: "A test package description",
					},
					nuget.SearchResult{
						ID:          "Some.Other.Package.Name",
						Version:     "2.0.10",
						Description: "A test package description",
					},
				},
			}
			statusCode = 200
		})
		It("should upload the package", func() {
			// //request := out.Request{
			// 	Source: nugetresource.Source{
			// 		NugetSource:  server.URL() + "/somefeed/api/v3/index.json",
			// 		NugetAPIKey:  "an_api_key",
			// 		NugetTimeout: 300,
			// 		PackageID:    "Some.Package.Name",
			// 	},
			// 	Params: out.Params{
			// 		PackagePath: "DotnetResource.TestApplication.1.0.28.nupkg",
			// 	},
			// }
			//_, _, _ := out.Execute(request, "../test_files")			
			//Ω(err).Should(HaveOccurred())
			//Ω(response).ShouldNot(BeNil())
		})
	})
})
