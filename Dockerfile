# syntax=docker/dockerfile:1
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY .dist/domovyk /app/domovyk
COPY map/domovyk-map.yaml /app/map/domovyk-map.yaml
USER nonroot:nonroot
ENTRYPOINT ["/app/domovyk"]
