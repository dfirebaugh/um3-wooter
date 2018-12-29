FROM golang:1.8.5-jessie
# install glide
RUN go get github.com/joho/godotenv
# create a working directory
WORKDIR /go/src/app
# add source code
ADD . src
# add env file
ADD .env /go/src/app
# build main.go
RUN go build src/main.go src/time.go
# run the binary
CMD ["./main"]