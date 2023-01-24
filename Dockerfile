# Build Step
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git tzdata g++ make ca-certificates


# Source https://chemidy.medium.com/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
# See https://stackoverflow.com/a/55757473/12429735RUN 
ENV USER=appuser
ENV UID=10001 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"



WORKDIR /code

COPY go.* .
RUN go mod download
ADD . /code/

# https://awstip.com/containerize-go-sqlite-with-docker-6d7fbecd14f0
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o indiego -a .

RUN mkdir -p /code/db

FROM alpine

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /app
COPY --from=builder /code/indiego /app/indiego

COPY config.json /app/config/config.json

EXPOSE 8080

ENTRYPOINT ["/app/indiego", "-config", "/app/config/config.json"]