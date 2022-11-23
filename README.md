# IndieGo

This project is a simple, modular but extendable webserver to enable your website to join the [IndieWeb](https://indieweb.org).

IndieGo currently supports native comments, sending and receiving [webmentions](https://indieweb.org/Webmention), [IndieAuth](https://indieweb.org/IndieAuth) as well as [Micropub](https://indieweb.org/Micropub). All components are built in a modular way, they can be extended for your own needs or disabled by just commenting them out in the [main method](/main.go). 

The API can be queried for comments and webmentions of all pages or just of a single page. 

> [Blogposts about IndieGo](https://tiim.ch/tags/indiego)

![Image of the Go gopher with a speech bubble](/go-comment-api-image.svg)

## Demo

You can try out this project on my [Blogpost about this project](https://tiim.ch/blog/2022-07-12-first-go-project-commenting-api).

## Installation

### Using docker compose

The easiest way to run the go-comment-api is via docker compose. There is a [sample docker-compose.yml](/docker-compose.yml) that I use to host the comments for my website. To see how I deploy it on my webserver see [deploy.sh](/deploy.sh).

### Compile it to a static binary
You need a recent go version installed. Run the following command to compile:

```sh
go mod download
CGO_ENABLED=0 go build -o comment-api -a .
```

This binary contains all needed data and templates. Run it with 
```sh
./comment-api
```

## Development

### Running Tests

```sh
go test ./...
```

### Running migrations in the cli

```sh
# Up
goose -dir model/sqlite-migrations/ sqlite3 db/comments.sqlite up
# Down
goose -dir model/sqlite-migrations/ sqlite3 db/comments.sqlite down
```

### Tools

#### Exposing tunnel to the web

```
npx ngrok http 8080
```

For testing IndieAuth:

```
INDIE_CANONICAL_URL=https://<ngok-url> go run .
```

#### Testing indieauth websub / indie reader with aperture

- Expose app to web via ngrok
- Register with ngrok url
- Restart app with `APPERTURE_ID=xxx`


#### Node Webmention Testpinger

[Source](https://github.com/voxpelli/node-webmention-testpinger)

- `npx webmention-testpinger --endpoint=http://localhost:8080/wm/webmentions --target https://tiim.ch/target -p 8081`
- `npx webmention-testpinger --endpoint=http://localhost:8080/wm/webmentions --target http://localhost:5173/projects/lenex-split-sheet -p 8081`

#### Sending single webmention with curl

```sh
cd test-data/
python3 -m http.server
curl -i -d source=http://localhost:8000/html/webmention-rocks.html -d target=https://tiim.ch/blog/2022-07-12-first-go-project-commenting-api http://localhost:8080/wm/webmentions
```

#### Testing sending webmentions

```sh
cd test-data/
python3 -m http.server
```