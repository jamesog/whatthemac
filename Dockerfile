FROM golang:1.14 AS build

WORKDIR /go/src/github.com/jamesog/mac.jog.li
COPY . .

RUN go get -d -v ./... && \
	CGO_ENABLED=0 go install -a -tags netgo -ldflags '-w' -v ./cmd/mac.jog.li && \
	echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

FROM scratch
COPY --from=build /go/bin/mac.jog.li /mac.jog.li
COPY --from=build /etc_passwd /etc/passwd

USER nobody

CMD ["/mac.jog.li"]
