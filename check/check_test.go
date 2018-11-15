package check_test

import (
	
	"github.com/onsi/gomega/ghttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/miclip/nuget-resource/check"
	"github.com/miclip/nuget-resource"
	"github.com/miclip/nuget-resource/nuget"
)

var _ = Describe("GetLatestVersion", func() {
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
						ID:      server.URL() + "/somefeed/api/v3/query",
						Type:    "SearchQueryService/3.0.0-beta",
						Comment: "Query endpoint of NuGet Search service.",
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
						Version:     "2.0.10",
						Description: "A test package description",
						Versions: []nuget.Version{
							nuget.Version{
								ID: server.URL() + "/somefeed/api/v3/registration1/dotnetresource.testlibraryone/1.0.1.json",
								Version: "2.0.9",
								Downloads: 3,
							},
							nuget.Version{
								ID: server.URL() + "/somefeed/api/v3/registration1/dotnetresource.testlibraryone/1.0.1.json",
								Version: "2.0.10",
								Downloads: 4,
							},
						},
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

		It("returns the versions", func() {
			request := check.Request{
				Source: nugetresource.Source{
					NugetSource: server.URL() + "/somefeed/api/v3/index.json",
					PackageID: "Some.Package.Name",
					PreRelease: false,
				},
				Version: nugetresource.Version{},
			}
			response, err := check.Execute(request)
			Ω(err).Should(Succeed())
			Ω(len(response)).Should(Equal(2))
			Ω(response[0].Version).Should(Equal("2.0.9"))
			Ω(response[0].PackageID).Should(Equal("Some.Package.Name"))
			Ω(response[1].Version).Should(Equal("2.0.10"))
			Ω(response[1].PackageID).Should(Equal("Some.Package.Name"))
		})

	
	})
})
