# バイナリにするとテンプレートを読み込めなかったので、コメントアウト
# FROM golang:1.15-alpine AS build

# WORKDIR /workspace
# COPY go.sum go.mod  ./
# RUN go mod download
# COPY ./ ./

# RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /workspace/go_gin_practice ./cmd/main.go

# FROM alpine as final

# WORKDIR /

# COPY --from=build /workspace/go_gin_practice /

# EXPOSE 8080

# ENTRYPOINT [ "./go_gin_practice" ]

# RUN apk --no-cache add tzdata && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# 以下でそのままイメージに入れる
FROM golang:latest
RUN mkdir /go/src/app
WORKDIR /go/src/app
ADD . /go/src/app
VOLUME /go/src/app
RUN go mod download
EXPOSE 8080
CMD "go" "run" "cmd/main.go"