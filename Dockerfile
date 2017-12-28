FROM golang:alpine AS build
WORKDIR /go/src/github.com/cuigh/xdoc/
ADD . .
RUN dep ensure
RUN go build -ldflags "-s -w"

FROM alpine:3.6
LABEL maintainer="noname@live.com"
ENV XDOC_DIR /docs
WORKDIR /app
COPY --from=build /go/src/github.com/cuigh/xdoc/xdoc .
COPY --from=build /go/src/github.com/cuigh/xdoc/config ./config/
COPY --from=build /go/src/github.com/cuigh/xdoc/assets ./assets/
COPY --from=build /go/src/github.com/cuigh/xdoc/views ./views/
EXPOSE 8000
ENTRYPOINT ["/app/xdoc"]
