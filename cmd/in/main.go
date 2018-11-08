package main

import (
	"log"
	"os"
	"encoding/json"
	"github.com/miclip/nuget-resource"
	"github.com/miclip/nuget-resource/in"
)

func main() {

	if len(os.Args) < 2 {
		nugetresource.Sayf("usage: %s <sources directory>\n", os.Args[0])
		os.Exit(1)
	}

	var request in.Request
	inputRequest(&request)

	response, output, err := in.Execute(request, os.Args[1])
	nugetresource.Sayf(string(output))

	if err != nil {
		log.Fatal(err)
	}

	outputResponse(response)
}

func inputRequest(request *in.Request) {
	if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
		nugetresource.Fatal("reading request from stdin", err)
	}
}

func outputResponse(response in.Response) {
	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		nugetresource.Fatal("writing response to stdout", err)
	}
}