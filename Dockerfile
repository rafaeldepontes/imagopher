FROM golang:1.25-alpine AS build

WORKDIR /cmd

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /app .

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=build /app /app

USER nonroot:nonroot

ENTRYPOINT [ "/app" ]