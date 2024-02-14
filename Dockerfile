FROM golang:latest

WORKDIR /go

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["go" ,"run" ,"."]