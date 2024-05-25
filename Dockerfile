FROM golang:1.22.3-alpine

# set the working dir
WORKDIR /app

# copy files to the working dir
COPY go.mod go.sum ./

# download and install dependencies
RUN go mod download

# copy entire sourfce code to working dif
COPY . .

# build the app
RUN go build -o main .

# expose the port
EXPOSE 5000

# run the app
CMD ["./main"]
