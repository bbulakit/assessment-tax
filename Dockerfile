FROM golang:1.21.6-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test --tags=integration -v ./...

RUN go build -o ./out/assessment-tax .

FROM alpine:3.16.2
COPY --from=build-base /app/out/assessment-tax /app/assessment-tax

ENV PORT="8080"
ENV DATABASE_URL="host=host.docker.internal port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable"
ENV ADMIN_USERNAME="adminTax"
ENV ADMIN_PASSWORD="admin!"

EXPOSE $PORT

CMD ["/app/assessment-tax"]

#Test commands
#docker run -p 8080:8029 -e PORT=8029 -e ADMIN_USERNAME=adminTax2 -e ADMIN_PASSWORD=P@ssw0rd -e DATABASE_URL="host=172.17.68.49 port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable" assessment-tax