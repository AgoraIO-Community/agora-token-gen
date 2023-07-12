# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine
RUN apk add git ca-certificates --update

# fetch dependencies github
# RUN go get -u github.com/gin-gonic/gin

ADD . /go/src/github.com/AgoraIO-Community/agora-token-gen

# # fetch dependencies from github (Gin and Agora Token Service)
# RUN go install github.com/gin-gonic/gin@latest
# # RUN go install github.com/AgoraIO-Community/agora-token-gen
# ADD . /go/src/github.com/AgoraIO-Community/agora-token-gen

ARG SERVER_PORT
ENV SERVER_PORT $SERVER_PORT

# move to the working directory
WORKDIR $GOPATH/src/github.com/AgoraIO-Community/agora-token-gen
# Build the token server command inside the container.
RUN go build -o agora-token-gen -v cmd/main.go
# RUN go run main.go
# Run the token server by default when the container starts.
ENTRYPOINT ./agora-token-gen

# Document that the service listens on port 8080.
EXPOSE 8080