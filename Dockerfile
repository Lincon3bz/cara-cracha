FROM golang:1.26 AS compilacao
WORKDIR /src
COPY go.mod go.sum ./
RUN go.mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /cara-cracha ./cmd/web

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=compilacao /cara-cracha /cara-cracha
EXPOSE 8080
ENTRYPOINT ["/cara-cracha"]