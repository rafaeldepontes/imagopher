FROM golang:1.26-alpine AS build

WORKDIR /cmd

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /app .

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=build /app /app
COPY --chown=nonroot:nonroot ./private /private

USER nonroot:nonroot

ENTRYPOINT [ "/app" ]