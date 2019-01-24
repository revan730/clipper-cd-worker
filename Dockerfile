FROM golang:latest AS BUILDER

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
ADD . .
RUN CGO_ENABLED=0 GO111MODULE=on go build -o clipper-cd

FROM wernight/kubectl:1.6.4

COPY --from=builder /go/src/app/clipper-cd /bin/clipper-cd

EXPOSE 8080
ENTRYPOINT ["/bin/clipper-cd"]
CMD ["start"]
