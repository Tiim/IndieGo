# Build Step
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git tzdata g++ make


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
    CGO_ENABLED=1 GOOS=linux \
    go build -o comment-api -a -ldflags '-linkmode external -extldflags "-static"' .

USER appuser:appuser
RUN mkdir -p /code/db

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

#USER appuser:appuser
WORKDIR /app
COPY --from=builder /code/comment-api /app/comment-api
COPY --from=builder /code/db /app/db

EXPOSE 8080

CMD ["/app/comment-api"]