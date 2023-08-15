FROM golang as base
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/ldap-api
FROM golang:1.16.11-bullseye
#ENV apitoken==Enter-Token-For-API-Authentication
#ENV LDAP_SERVER=Main-AD
#ENV BIND_ADDRESS=:8090

RUN echo "The ENV variable value is $apitoken"
WORKDIR /app
COPY --from=base /app/ldap-api .
COPY --from=base /app/swagger.yaml .
EXPOSE 8090
ENTRYPOINT ["/bin/sh", "-c" , "echo Main-AD-IP Main-AD-Domain >> /etc/hosts && /app/ldap-api && echo Main-AD2-IP Main-AD2-Domain >> /etc/hosts && /app/ldap-api"]
