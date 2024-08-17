[![progress-banner](https://backend.codecrafters.io/progress/http-server/6b4337b9-e9e7-4c47-8d17-3d5fe6f02261)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

[HTTP](https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol) is the
protocol that powers the web. This repo builds a HTTP/1.1 server
that is capable of serving multiple clients.

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to try the challenge.

# Running the app

`go run app/server.go --directory /tmp`

## Sample curls to try

```sh
curl -v http://localhost:4221/

curl -v http://localhost:4221/user-agent -H "User-Agent: grape/grape-mango"

curl -v http://localhost:4221/echo/blueberry

curl -v -X POST http://localhost:4221/files/orange_orange_strawberry_raspberry -H "Content-Length: 56" -H "Content-Type: application/octet-stream" -d 'apple apple pear strawberry apple mango blueberry orange'

curl -v -H "Accept-Encoding: gzip" http://localhost:4221/echo/abc | gzip -d
```
