FROM golang:1.14

WORKDIR /app
COPY app.go .

# Install the dependencies
RUN go get -d github.com/gorilla/mux
RUN go get -d go.mongodb.org/mongo-driver/mongo
RUN go get -d go.mongodb.org/mongo-driver/mongo/options


CMD ["go", "run", "app.go"]