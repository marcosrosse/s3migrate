# Build CLI app binary
FROM golang:1.21.0 AS build-stage

WORKDIR /app

COPY . ./ 

RUN go mod download

# Build the binary file
RUN CGO_ENABLED=0 GOOS=linux go build -o /s3migrate cmd/main.go

# Deploy the CLI app binary into a lean image
FROM alpine AS release-stage

ENV S3_ENDPOINT=minio:9000 \
    S3_ACCESS_KEY_ID=minioadmin \
    S3_SECRET_ACCESS_KEY=minioadmin \
    S3_USE_SSL=false \
    POSTGRES_HOST=postgres \
    POSTGRES_PORT=5432 \
    POSTGRES_USERNAME=postgres \
    POSTGRES_PASSWORD=postgres \
    POSTGRES_DBNAME=proddatabase \
    S3_SRC_BUCKET="legacy-s3" \
    S3_DST_BUCKET="production-s3" \
    S3_SRC_PATH_OBJ="image" \
    S3_DST_PATH_OBJ="avatar"

WORKDIR /

COPY --from=build-stage /s3migrate /s3migrate

ENTRYPOINT [ "/s3migrate" ]