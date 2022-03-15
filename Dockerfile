FROM golang:1.17-alpine AS build

WORKDIR /wasp-simple-h264-payload-parser
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/wasp-simple-h264-payload-parser

FROM alpine
COPY --from=build /bin/wasp-simple-h264-payload-parser /bin/wasp-simple-h264-payload-parser
CMD ["/bin/wasp-simple-h264-payload-parser"]