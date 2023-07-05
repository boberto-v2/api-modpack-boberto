FROM golang:latest
WORKDIR /app
COPY . .

RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
RUN go install github.com/swaggo/swag/cmd/swag@latest
ENTRYPOINT CompileDaemon -command="go run main.go" -exclude-dir=.git -polling
EXPOSE 80