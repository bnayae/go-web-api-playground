FROM golang:1.11.5

# Install mux
RUN go get github.com/gorilla/mux

# Create app directory
WORKDIR /usr/go/server
Copy ./src/ .
# Expose the application on port 8080
EXPOSE 80

Run ls

CMD ["go", "run", "main.go"]