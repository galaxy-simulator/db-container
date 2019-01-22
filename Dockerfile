FROM golang:latest

WORKDIR /home

COPY *.go /home/

RUN ["mkdir", "/exports"]
RUN ["go", "get", "github.com/gorilla/mux"]
RUN ["go", "get", "git.darknebu.la/GalaxySimulator/structs"]
RUN ["ls", "-l"]

ENTRYPOINT ["go", "run", "."]
