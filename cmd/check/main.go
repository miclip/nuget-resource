package main

import (
	"log"
	"os"
	"encoding/json"

	"github.com/miclip/nuget-resource/check"
	"github.com/miclip/nuget-resource"
)

func main() {
	var request check.Request
	inputRequest(&request)

	response, err := check.Execute(request)
	if err != nil {
		log.Fatal(err)
	}		
	outputResponse(response)
}

func inputRequest(request *check.Request) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		nugetresource.Fatal("reading request from stdin", err)
	}
}

func outputResponse(response check.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		nugetresource.Fatal("writing response to stdout", err)
	}
}