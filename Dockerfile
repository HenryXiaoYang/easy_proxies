# Stage 1: Build simple-obfs plugin (obfs-local binary)
FROM debian:bookworm-slim AS obfs-builder
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    autoconf \
    libtool \
    libssl-dev \
    libpcre3-dev \
    libev-dev \
    automake \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ARG OBFS_VERSION=0.0.5
RUN git clone https://github.com/shadowsocks/simple-obfs.git /tmp/simple-obfs && \
    cd /tmp/simple-obfs && \
    git checkout v${OBFS_VERSION} && \
    git submodule update --init --recursive && \
    ./autogen.sh && \
    ./configure --disable-documentation && \
    make CFLAGS="-Wno-error=stringop-overread" && \
    make install

# Stage 2: Build Go application
FROM --platform=$BUILDPLATFORM golang:1.24 AS builder
ARG TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
ARG GOPROXY=https://proxy.golang.org,direct
RUN go env -w GOPROXY=${GOPROXY} && go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -tags "with_utls with_quic with_grpc with_wireguard with_gvisor" -o easy-proxies ./cmd/easy_proxies

FROM debian:bookworm-slim AS runtime
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        libev4 \
        libpcre3 \
    && rm -rf /var/lib/apt/lists/* \
    && useradd -r -u 10001 easy \
    && mkdir -p /etc/easy-proxies \
    && chown -R easy:easy /etc/easy-proxies
WORKDIR /app
COPY --from=obfs-builder /usr/local/bin/obfs-local /usr/local/bin/obfs-local
COPY --from=builder /src/easy-proxies /usr/local/bin/easy-proxies
COPY --chown=easy:easy config.example.yaml /etc/easy-proxies/config.yaml
# Pool/Hybrid mode: 2323, Management: 9091, Multi-port/Hybrid mode: 24000-24200
EXPOSE 2323 9091 24000-24200
USER easy
ENTRYPOINT ["/usr/local/bin/easy-proxies"]
CMD ["--config", "/etc/easy-proxies/config.yaml"]
