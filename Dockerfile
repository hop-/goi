FROM golang:1.22-alpine AS dev-build
RUN apk add --no-cache make
WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY Makefile .
RUN make setup
COPY . .

FROM dev-build AS dev
COPY --from=main / /
ENV HOST_ENV=development
ENV HOST_CONFIG_DIR=configs
CMD ["go", "run", "./cmd/goi/main.go"]

FROM dev-build AS prod-build
RUN make

FROM alpine:3.20 AS prod
ENV HOST_CONFIG_DIR=configs
WORKDIR /etc/goi
COPY --from=prod-build /build/configs ./configs
COPY --from=prod-build /build/bin/goi /usr/local/bin
ENV HOST_ENV=production
CMD ["goi"]
