#!/usr/bin/env sh

go run honnef.co/go/tools/cmd/staticcheck ./...
go run github.com/onsi/ginkgo/ginkgo -r
