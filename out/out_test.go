package out_test

import (
	"github.com/miclip/nuget-resource"
	"github.com/miclip/nuget-resource/nuget"
	"github.com/miclip/nuget-resource/out"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("out", func() {


	It("should output an empty JSON list", func() {
		request := out.Request{
			Source: nugetresource.Source{
				Framework: "netcoreapp2.1",
				Runtime:   "ubuntu.14.04-x64",
			},
			Params: out.Params{
				Project:    "/path/project.csproj",
				TestFilter: "A_Filter",
			},
		}
		response, _, err := out.Execute(request, "/tmp")
		Ω(response).ShouldNot(BeNil())
		Ω(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("generateNextVersion", func() {
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

		// It("returns a version with patch incremented", func() {
		// 	request := out.Request{
		// 		Source: nugetresource.Source{
		// 			Framework:   "netcoreapp2.1",
		// 			Runtime:     "ubuntu.14.04-x64",
		// 			NugetSource: server.URL() + "/somefeed/api/v3/index.json",
		// 		},
		// 		Params: out.Params{
		// 			Project:     "/path/project.csproj",
		// 			TestFilter:  "A_Filter",
		// 			Version:     "1.0.*",			
		// 		},
		// 	}
		// 	version, err := out.GenerateNextVersion(request, "Some.Package.Name")
		// 	Ω(err).Should(Succeed())
		// 	Ω(version).Should(Equal("2.0.3"))
		// })

		// It("returns a version with minor & patch incremented", func() {
		// 	request := out.Request{
		// 		Source: nugetresource.Source{
		// 			Framework:   "netcoreapp2.1",
		// 			Runtime:     "ubuntu.14.04-x64",
		// 			NugetSource: server.URL() + "/somefeed/api/v3/index.json",
		// 		},
		// 		Params: out.Params{
		// 			Project:     "/path/project.csproj",
		// 			TestFilter:  "A_Filter",
		// 			Version:     "1.*.*",
		// 		},
		// 	}
		// 	version, err := out.GenerateNextVersion(request, "Some.Package.Name")
		// 	Ω(err).Should(Succeed())
		// 	Ω(version).Should(Equal("2.1.3"))
		// })

		// It("returns a version with minor incremented", func() {
		// 	request := out.Request{
		// 		Source: nugetresource.Source{
		// 			Framework:   "netcoreapp2.1",
		// 			Runtime:     "ubuntu.14.04-x64",
		// 			NugetSource: server.URL() + "/somefeed/api/v3/index.json",
		// 		},
		// 		Params: out.Params{
		// 			Project:     "/path/project.csproj",
		// 			TestFilter:  "A_Filter",
		// 			Version:     "1.*.0",
		// 		},
		// 	}
		// 	version, err := out.GenerateNextVersion(request, "Some.Package.Name")
		// 	Ω(err).Should(Succeed())
		// 	Ω(version).Should(Equal("2.1.2"))
		// })

		// It("returns a version with major, minor & patch incremented", func() {
		// 	request := out.Request{
		// 		Source: nugetresource.Source{
		// 			Framework:   "netcoreapp2.1",
		// 			Runtime:     "ubuntu.14.04-x64",
		// 			NugetSource: server.URL() + "/somefeed/api/v3/index.json",
		// 		},
		// 		Params: out.Params{
		// 			Project:     "/path/project.csproj",
		// 			TestFilter:  "A_Filter",
		// 			Version:     "*.*.*",
		// 		},
		// 	}
		// 	version, err := out.GenerateNextVersion(request, "Some.Package.Name")
		// 	Ω(err).Should(Succeed())
		// 	Ω(version).Should(Equal("3.1.3"))
		// })

		// It("returns error as version semantics don't match", func() {
		// 	request := out.Request{
		// 		Source: nugetresource.Source{
		// 			Framework:   "netcoreapp2.1",
		// 			Runtime:     "ubuntu.14.04-x64",
		// 			NugetSource: server.URL() + "/somefeed/api/v3/index.json",
		// 		},
		// 		Params: out.Params{
		// 			Project:     "/path/project.csproj",
		// 			TestFilter:  "A_Filter",
		// 			Version:     "1.0.1-dev.*",
		// 		},
		// 	}
		// 	_,err := out.GenerateNextVersion(request, "Some.Package.Name")
		// 	Ω(err).ShouldNot(Succeed())
		// })
	})
})
