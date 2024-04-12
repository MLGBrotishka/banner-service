FROM golang:1.22.1
WORKDIR /app

COPY . . 

RUN make tidy

RUN make build

EXPOSE 8000

CMD ["./myapp"]