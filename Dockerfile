
FROM golang:1.18


WORKDIR /ufils

COPY . .

EXPOSE 8080

CMD ["go","run","."]
