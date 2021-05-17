# golang http-server demo

```bash

docker run -d --restart=always --name http-server -v $PWD:/app/static -w /app -p 38080:8080 10.58.10.201:55000/httpserver

```