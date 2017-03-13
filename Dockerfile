FROM golang:alpine
RUN apk add --no-cache iptables ca-certificates git && \
	update-ca-certificates

ENV DO_KEY ""
ENV DO_TAG ""

RUN ln -s go/src/github.com/tam7t/droplan/ /droplan
COPY ./ /go/src/github.com/tam7t/droplan/

WORKDIR /droplan
RUN go get -u github.com/kardianos/govendor
RUN govendor build

ENTRYPOINT ["/droplan/docker-run.sh"]
