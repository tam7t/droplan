FROM alpine:3.4
RUN apk add --no-cache iptables ca-certificates && \
	update-ca-certificates

ENV DO_KEY ""
ENV DO_TAG ""

ADD droplan /droplan
ADD docker-run.sh /docker-run.sh

CMD ["/docker-run.sh"]
