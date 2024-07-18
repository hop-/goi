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
# RUN sysctl -w net.core.rmem_max=7500000 \
#     && sysctl -w net.core.wmem_max=7500000
CMD ["go", "run", "./cmd/goi/main.go"]

FROM dev-build AS prod-build
RUN make

FROM alpine:3.20 AS prod
ENV HOST_ENV=production
ENV HOST_CONFIG_DIR=configs
# RUN sysctl -w net.core.rmem_max=7500000 \
#     && sysctl -w net.core.wmem_max=7500000
WORKDIR /etc/goi
COPY --from=prod-build /build/configs ./configs
COPY --from=prod-build /build/bin/goi /usr/local/bin
CMD ["goi"]
