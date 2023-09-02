FROM golang:1.20-alpine
RUN apk add --update npm

WORKDIR /app

COPY ./ ./

RUN cd web && npm install && npm run build
RUN go mod download
RUN go build -o /sun-forecast

CMD [ "/sun-forecast" ]