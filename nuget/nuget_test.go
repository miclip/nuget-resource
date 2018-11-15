package nuget_test

import (
	"net/http"
	"context"
	"encoding/xml"

	"github.com/miclip/nuget-resource/nuget"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("SearchQueryService", func() {
	var server *ghttp.Server
	var returnedSearchResults nuget.SearchResults
	var statusCode int
	var client nuget.NugetClientv3

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = nuget.NewNugetClientv3(server.URL() + "/somefeed/api/v3/index.json")
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/somefeed/api/v3/query"),
				ghttp.RespondWithJSONEncodedPtr(&statusCode, &returnedSearchResults),
			),
		)
	})

	AfterEach(func() {
		server.Close()
	})

	Context("When the server returns a Search Result", func() {
		BeforeEach(func() {

			returnedSearchResults = nuget.SearchResults{
				TotalHits:  1000,
				Index:      "index",
				LastReopen: "2018-10-22T22:45:00.425508Z",
				Data: []nuget.SearchResult{
					nuget.SearchResult{
						ID:          "Some.Package.Name",
						Version:     "2.0.1",
						Description: "A test package description",
					},
				},
			}
			statusCode = 200
		})

		It("returns Service Index of a nuget feed", func() {
			r, err := client.SearchQueryService(context.Background(), server.URL()+"/somefeed/api/v3/query", "Some.Package.Name", true)
			Ω(err).Should(Succeed())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
			Ω(r).ShouldNot(BeNil())

		})
	})

	Context("when the server returns 500", func() {
		BeforeEach(func() {
			statusCode = 500
		})
		It("errors", func() {
			r, err := client.SearchQueryService(context.Background(), server.URL()+"/somefeed/api/v3/query", "Some.Package.Name", true)
			Ω(err).To(HaveOccurred())
			Ω(r).To(BeNil())

		})
	})

	Context("when the server returns 503", func() {
		BeforeEach(func() {
			statusCode = 503
		})

		It("errors", func() {
			r, err := client.SearchQueryService(context.Background(), server.URL()+"/somefeed/api/v3/query", "Some.Package.Name", true)
			Ω(err).To(HaveOccurred())
			Ω(r).To(BeNil())
		})
	})
})

var _ = Describe("ServiceIndex", func() {
	var server *ghttp.Server
	var returnedServiceIndex nuget.ServiceIndex
	var statusCode int
	var client nuget.NugetClientv3

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = nuget.NewNugetClientv3(server.URL() + "/somefeed/api/v3/index.json")
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/somefeed/api/v3/index.json"),
				ghttp.RespondWithJSONEncodedPtr(&statusCode, &returnedServiceIndex),
			),
		)
	})

	AfterEach(func() {
		server.Close()
	})

	Context("When when the server returns a Service Index", func() {
		BeforeEach(func() {
			returnedServiceIndex = nuget.ServiceIndex{
				Version: "3.0.0",
				Resources: []nuget.Resource{
					nuget.Resource{
						ID:      "https://www.nuget.org/somefeed/api/v3/query",
						Type:    "SearchQueryService",
						Comment: "Query endpoint of NuGet Search service.",
					},
				},
			}
			statusCode = 200
		})

		It("returns Service Index of a nuget feed", func() {
			r, err := client.GetServiceIndex(context.Background())
			Ω(err).Should(Succeed())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
			Ω(r).ShouldNot(BeNil())
			Ω(r.Resources).Should(HaveLen(1))
			Ω(r.Resources[0].Type).To(Equal("SearchQueryService"))
			Ω(r.Version).To(Equal("3.0.0"))
		})
	})

	Context("when the server returns 500", func() {
		BeforeEach(func() {
			statusCode = 500
		})

		It("errors", func() {
			r, err := client.GetServiceIndex(context.Background())
			Ω(err).To(HaveOccurred())
			Ω(r).To(BeNil())

		})
	})

	Context("when the server returns 503", func() {
		BeforeEach(func() {
			statusCode = 503
		})

		It("errors", func() {
			r, err := client.GetServiceIndex(context.Background())
			Ω(err).To(HaveOccurred())
			Ω(r).To(BeNil())
		})
	})
})

var _ = Describe("GetPackageVersions", func() {
	var server *ghttp.Server
	var returnedServiceIndex nuget.ServiceIndex
	var returnedSearchResults nuget.SearchResults
	var statusCode int
	var client nuget.NugetClientv3

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = nuget.NewNugetClientv3(server.URL() + "/somefeed/api/v3/index.json")
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

	Context("When the client returns a Version", func() {
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
						Version:     "2.0.1",
						Description: "A test package description",
						Versions: []nuget.Version{
							nuget.Version{
								ID: server.URL() + "/somefeed/api/v3/registration1/dotnetresource.testlibraryone/1.0.1.json",
								Version: "2.0.0",
								Downloads: 3,
							},
							nuget.Version{
								ID: server.URL() + "/somefeed/api/v3/registration1/dotnetresource.testlibraryone/1.0.1.json",
								Version: "2.0.1",
								Downloads: 3,
							},
						},
					},
					nuget.SearchResult{
						ID:          "Some.Other.Package.Name",
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
								Downloads: 3,
							},
						},
					},
				},
			}
			statusCode = 200
		})

		It("returns all versions for a particular package", func() {
			r, err := client.GetPackageVersions(context.Background(), "Some.Package.Name", false)
			Ω(err).Should(Succeed())
			Ω(server.ReceivedRequests()).Should(HaveLen(2))
			Ω(r).ShouldNot(BeNil())
			Ω(r).Should(HaveLen(2))
			Ω(r[0].Version).To(Equal("2.0.0"))
			Ω(r[1].Version).To(Equal("2.0.1"))
		})
	})
})

var _ = Describe("Versioning", func(){
	Context("AutoIncrementVersions", func() {
		It("increments patch", func() {
			client := nuget.NewNugetClientv3("http://nuget.org/somefeed/api/v3/index.json")
			version, err := client.AutoIncrementVersion("1.0.*","1.1.9") 
			Ω(err).Should(Succeed())
			Ω(version).Should(Equal("1.1.10"))
		})
		It("increments minor", func() {
			client := nuget.NewNugetClientv3("http://nuget.org/somefeed/api/v3/index.json")
			version, err := client.AutoIncrementVersion("1.*.0","1.1.9") 
			Ω(err).Should(Succeed())
			Ω(version).Should(Equal("1.2.9"))
		})
		It("increments major", func() {
			client := nuget.NewNugetClientv3("http://nuget.org/somefeed/api/v3/index.json")
			version, err := client.AutoIncrementVersion("*.1.0","1.1.9") 
			Ω(err).Should(Succeed())
			Ω(version).Should(Equal("2.1.9"))
		})
		It("increments all", func() {
			client := nuget.NewNugetClientv3("http://nuget.org/somefeed/api/v3/index.json")
			version, err := client.AutoIncrementVersion("*.*.*","1.1.9") 
			Ω(err).Should(Succeed())
			Ω(version).Should(Equal("2.2.10"))
		})
		It("increments build with suffix", func() {
			client := nuget.NewNugetClientv3("http://nuget.org/somefeed/api/v3/index.json")
			version, err := client.AutoIncrementVersion("1.0.0-dev.*","1.0.0-dev.187") 
			Ω(err).Should(Succeed())
			Ω(version).Should(Equal("1.0.0-dev.188"))
		})
		It("errors when schemas don't match", func() {
			client := nuget.NewNugetClientv3("http://nuget.org/somefeed/api/v3/index.json")
			version, err := client.AutoIncrementVersion("1.0.0-dev.*","1.0.0") 
			Ω(err).ShouldNot(Succeed())
			Ω(version).Should(Equal(""))
		})
		It("returns unmodified version if no mask", func() {
			client := nuget.NewNugetClientv3("http://nuget.org/somefeed/api/v3/index.json")
			version, err := client.AutoIncrementVersion("1.0.0","1.1.9") 
			Ω(err).Should(Succeed())
			Ω(version).Should(Equal("1.1.9"))
		})
		It("returns unmodified version if no mask", func() {
			client := nuget.NewNugetClientv3("http://nuget.org/somefeed/api/v3/index.json")
			version, err := client.AutoIncrementVersion("","1.1.9") 
			Ω(err).Should(Succeed())
			Ω(version).Should(Equal("1.1.9"))
		})
	})
})

var _ = Describe("CreateNuspec", func() {
	Context("nuspec encoding", func() {
		It("returns a valid nuspec ", func() {
			client := nuget.NewNugetClientv3("http://nuget.org/somefeed/api/v3/index.json")
			nuspec := client.CreateNuspec("packageID", "1.0.0", "Michael Lipscombe", "description", "owner")
			output, _ := xml.Marshal(nuspec)
			Ω(string(output)).To(Equal("<package xmlns=\"http://schemas.microsoft.com/packaging/2013/05/nuspec.xsd\"><metadata><id>packageID</id><version>1.0.0</version><authors>Michael Lipscombe</authors><owners>owner</owners><requireLicenseAcceptance>false</requireLicenseAcceptance><description>description</description></metadata></package>"))
		})
	})
	Context("nuspec from project file", func() {
		It("returns a valid nuspec based on a csproj with new version", func() {
			client := nuget.NewNugetClientv3("http://nuget.org/somefeed/api/v3/index.json")
			nuspec, err := client.CreateNuspecFromProject("../test_files/TestApplication.csproj", "9.9.9")
			Ω(err).Should(Succeed())
			Ω(nuspec.ID).Should(Equal("DotnetResource.TestApplication"))
			Ω(nuspec.Version).Should(Equal("9.9.9"))
			Ω(nuspec.Authors).Should(Equal("Michael Lipscombe"))
			Ω(nuspec.Owners).Should(Equal("Pivotal"))
			Ω(nuspec.Description).Should(Equal("A test application for dotnet-extensions"))
			Ω(nuspec.RequireLicenseAcceptance).Should(Equal(false))
		})
		It("returns a valid nuspec based on a csproj keeping default nuspec version", func() {
			client := nuget.NewNugetClientv3("http://nuget.org/somefeed/api/v3/index.json")
			nuspec, err := client.CreateNuspecFromProject("../test_files/TestApplication.csproj", "")
			Ω(err).Should(Succeed())
			Ω(nuspec.ID).Should(Equal("DotnetResource.TestApplication"))
			Ω(nuspec.Version).Should(Equal("1.0.0"))
			Ω(nuspec.Authors).Should(Equal("Michael Lipscombe"))
			Ω(nuspec.Owners).Should(Equal("Pivotal"))
			Ω(nuspec.Description).Should(Equal("A test application for dotnet-extensions"))
			Ω(nuspec.RequireLicenseAcceptance).Should(Equal(false))
		})
	})
})

var _ = Describe("UploadPackage", func() {
	var server *ghttp.Server
	var returnedServiceIndex nuget.ServiceIndex
	var statusCode int
	var client nuget.NugetClientv3

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = nuget.NewNugetClientv3(server.URL() + "/somefeed/api/v3/index.json")
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/somefeed/api/v3/index.json"),
				ghttp.RespondWithJSONEncodedPtr(&statusCode, &returnedServiceIndex),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/somefeed/api/v2/package"),
				ghttp.VerifyHeader(http.Header{"X-NuGet-ApiKey": []string{"an_api_key"}}),
			),
		)
	})

	AfterEach(func() {
		server.Close()
	})

	Context("when the package is uploaded", func() {
		BeforeEach(func() {
			returnedServiceIndex = nuget.ServiceIndex{
				Version: "3.0.0",
				Resources: []nuget.Resource{
					nuget.Resource{
						ID:      server.URL() + "/somefeed/api/v2/package",
						Type:    "PackagePublish/2.0.0",
						Comment: "Legacy gallery publish endpoint using the V2 protocol.",
					},
				},
			}			
			statusCode = 200
		})

		It("returns nil error", func() {
			err := client.PublishPackage(context.Background(), "an_api_key", "../test_files/DotnetResource.TestApplication.1.0.28.nupkg")
			Ω(err).Should(Succeed())
			Ω(server.ReceivedRequests()).Should(HaveLen(2))		
		})
	})
})

