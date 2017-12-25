FROM arm64v8/golang
WORKDIR /go/src/app
COPY . .

RUN go-wrapper download
RUN go-wrapper install 

CMD ["go-wrapper", "run"]