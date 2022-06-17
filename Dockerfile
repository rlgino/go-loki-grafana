FROM golang:alpine as build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o executable .

FROM alpine
WORKDIR /app
COPY --from=build /app/executable /app
RUN ls -la

ENTRYPOINT [ "/app/executable" ][]