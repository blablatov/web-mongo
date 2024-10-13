FROM golang:1.20
RUN wget https://github.com/blablatov/web-mongo.git
WORKDIR /web-mongo
COPY . .
RUN go test .
RUN go build -tags mongodb -o /web-mongo 
EXPOSE 8017
CMD ["./web-mongo"]
