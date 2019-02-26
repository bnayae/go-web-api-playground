# Dockerize

**docker build -t bnaya/go-server:v1 . **

docker run -it --rm --name bnaya-go-server -p 8088:80 -e GO_SERVER_PORT="80" bnaya/go-server:v1
