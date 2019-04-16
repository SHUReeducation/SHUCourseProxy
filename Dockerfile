FROM golang:1.12-alpine as builder
RUN apk update && apk --no-cache git add ca-certificates
ENV GO111MODULE=on
COPY . /go/src/SHUCourseProxy
WORKDIR /go/src/SHUCourseProxy
RUN go get && go build -o SHUCourseProxy .


FROM alpine
RUN apk update && apk --no-cache add ca-certificates
COPY --from=builder /go/src/SHUCourseProxy/SHUCourseProxy .
CMD ["./SHUCourseProxy"]
EXPOSE 8086

