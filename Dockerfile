FROM golang
WORKDIR /
RUN go get github.com/go-sql-driver/mysql && \
    go get github.com/gorilla/mux && \
    go get github.com/jmoiron/sqlx
ENV GOOS=linux
ENV GOARCH=amd64
COPY src/app.go /app.go
COPY src/app.conf /app.conf
ENTRYPOINT ["go", "run", "/app.go"]