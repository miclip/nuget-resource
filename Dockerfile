FROM golang:alpine as builder
COPY . /go/src/github.com/miclip/nuget-resource
ENV CGO_ENABLED 0
RUN go build -o /assets/in github.com/miclip/nuget-resource/cmd/in
RUN go build -o /assets/out github.com/miclip/nuget-resource/cmd/out
RUN go build -o /assets/check github.com/miclip/nuget-resource/cmd/check
RUN mkdir -p /tests/test_files
WORKDIR /go/src/github.com/miclip/nuget-resource
RUN cp ./test_files/* ../../../../../tests/test_files
RUN set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM alpine:edge AS resource
RUN apk add --no-cache bash tzdata ca-certificates unzip zip gzip tar
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*

ENV GOROOT=/goroot \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PATH=${PATH}:/goroot/bin:/gopath/bin
    

FROM resource AS tests
COPY --from=builder /tests /go-tests
COPY --from=builder /tests/test_files /test_files
WORKDIR /go-tests
RUN set -e; for test in /go-tests/*.test; do \
		$test; \
	done

FROM resource