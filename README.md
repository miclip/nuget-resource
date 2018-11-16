[![Build Status](https://travis-ci.org/miclip/nuget-resource.svg?branch=master)](https://travis-ci.org/miclip/nuget-resource)

# Nuget Resource

Pushes, downloads, extracts dotnet core libraries and applications to and from a nuget feed.

## Source Configuration

* `nuget_source`: *Required.* URL for nuget feed, only v3 API's are currently supported.

* `nuget_apikey`: *Required.* Nuget server API Key.

* `nuget_timeout`: *Optional.* The timeout for pushing packages to nuget, defaults to 300 seconds.

* `package_Id`: *Required for Check.* Package Name or PackageID as defined in the nuspec or csproj.

* `prerelease`: *Optional.* Whether the package is prerelease or not

## Behavior

### `check`: Detect a new version of a package

Requires the `package_id` and `prerelease` and nuget details in the source.

Resource:
```yml
resource_types:
  - name: nuget
    type: docker-image
    source:
      repository: miclip/nuget-resource
      tag: "latest"

resources:
  - name: nuget-get
    type: nuget 
    source:
      nuget_source: https://www.nuget.org/F/myfeed/api/v3/index.json
      nuget_apikey: {{nuget_apikey}}
      package_id: NugetResource.TestApplication
      prerelease: true
```
Job:
```yml
  name: deploy-service
  public: true
  serial: true
  plan:
  - get: nuget-get
    trigger: true
  - put: cf-resource
    params:
      manifest: nuget-get/manifest.yml
      path: nuget-get/
```

### `in`: Fetch a package from nuget.

Downloads and unpacks the package into the output of the resource.

### `out`: Push a package.

Given a package the resource will push it to the nuget feed. 

#### Parameters

* `package_path`: *Required.* Path to package file (nupkg)

## Example Configuration

### Resource Type

``` yaml
- name: nuget
  type: docker-image
  source:
    repository: miclip/nuget-resource
    tag: "latest"
```

### Resource

``` yaml
- name: nuget-out
  type: nuget 
  source:
    nuget_source: https://www.nuget.org/F/myfeed/api/v3/index.json
    nuget_apikey: {{nuget_apikey}}
    package_id: NugetResource.TestApplication
    prerelease: true
```

### Job

``` yaml
- name: build-service
  public: true
  serial: true
  plan:
  - get: app-repo
    trigger: true
  - task: build-and-pkg
    file: app-repo/ci/tasks/build.yml
  - put: nuget-out
    params: 
      package_path: build-output/*.nupkg
```

## Development

### Prerequisites

* golang is *required* - version 1.11.x is tested; earlier versions may also
  work.
* docker is *required* - version 18.06.x is tested; earlier versions may also
  work.
* dep is used for dependency management of the golang packages.

### Running the tests

The tests have been embedded with the `Dockerfile`; ensuring that the testing
environment is consistent across any `docker` enabled platform. When the docker
image builds, the test are run inside the docker container, on failure they
will stop the build.

Run the tests with the following command:

```sh
docker build -t nuget-resource .
```

### Examples 

Dotnet core MVC Application with tests:

https://github.com/miclip/nuget-resource-test-application

