FROM golang as builder
WORKDIR /go/src/github.com/llschuster/qasmoke/
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix .

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /go/src/github.com/llschuster/qasmoke/qasmoke /app/qasmoke
COPY --from=builder /go/src/github.com/llschuster/qasmoke/.env /app/.env
EXPOSE 5000
ENTRYPOINT ["./qasmoke"] 