FROM golang:1.22.1
WORKDIR /app

COPY . . 

RUN make tidy

CMD ["go", "test", "-v", "./..."]