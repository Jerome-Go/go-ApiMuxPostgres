# syntax=docker/dockerfile:1
#
# 1 image de base 
FROM golang:1.16-alpine
#
# 2 creation d'un repertoire à l'intérieur de l'image
WORKDIR /app
#
# 3 copie des fichiers go.mod et go.sum
COPY go.mod .
COPY go.sum .
#
# 4 installer les module go.mod et go.sum
RUN go mod download
#
# 5 copier tout les odres fichier go
COPY *.go ./
#
# 6 tout faire complitler
RUN go build -o /go-apipostgres
#
# 9 ouverture port 8000 sur le conteneur
EXPOSE 8000
#
# 10 commande pour executer notre app
CMD ["/go-apipostgres"]

