FROM golang:latest AS BUILDER

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
ADD . .
RUN CGO_ENABLED=0 GO111MODULE=on go build -o clipper-cd

FROM alpine:3.8

ADD https://storage.googleapis.com/kubernetes-release/release/v1.13.1/bin/linux/amd64/kubectl /usr/local/bin/kubectl

ENV HOME=/config

RUN set -x && \
    apk add --no-cache curl ca-certificates && \
    chmod +x /usr/local/bin/kubectl && \
    \
    # Create non-root user (with a randomly chosen UID/GUI).
    adduser kubectl -Du 2342 -h /config && \
    \
    # Basic check it works.
    kubectl version --client

USER kubectl

COPY --from=builder /go/src/app/clipper-cd /bin/clipper-cd

EXPOSE 8080
ENTRYPOINT ["/bin/clipper-cd"]
CMD ["start"]
