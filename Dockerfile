FROM golang:1.7-alpine

ENV DO_KEY ""
ADD droplan /droplan
ENTRYPOINT ["/droplan"]
