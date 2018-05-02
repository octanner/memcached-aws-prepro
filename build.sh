#!/bin/sh

cd /go/src
go get  "github.com/lib/pq"
go get  "github.com/aws/aws-sdk-go/aws"
go get  "github.com/aws/aws-sdk-go/aws/session"
go get  "github.com/aws/aws-sdk-go/service/elasticache"
go get  "github.com/nu7hatch/gouuid"

cd /go/src/oct-memcached-preprovision
go build oct-memcached-preprovision.go

