# This docker file is to be used without installing air
FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod download
# RUN go get github.com/githubnemo/CompileDaemon
# RUN go install github.com/githubnemo/CompileDaemon
# RUN go install github.com/swaggo/swag/cmd/swag@latest
# ENTRYPOINT CompileDaemon -command="go run main.go" -exclude-dir=.git -polling

RUN CGO_ENABLED=0 GOOS=linux go build -o /apimodpackboberto

CMD ["/apimodpackboberto"]
EXPOSE 80