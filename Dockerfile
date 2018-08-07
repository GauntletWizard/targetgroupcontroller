FROM golang:stretch

RUN apt-get update && apt-get install git
ADD . /go/src/github.com/lifeonair/targetgroupcontroller
RUN cd /go/src/github.com/lifeonair/targetgroupcontroller && go get && go build

FROM debian:stretch

COPY --from=0 /go/bin/targetgroupcontroller /bin/targetgroupcontroller
RUN apt-get update && apt-get install -y ca-certificates
ENTRYPOINT ["/bin/targetgroupcontroller"]
