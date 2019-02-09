#!/usr/bin/env bash
curl -X POST http://localhost:8081/deleteStars
curl -X POST http://localhost:8081/deleteNodes 
curl -X POST --data "w=1000" http://localhost:8081/new
curl -X POST --data "x=300&y=300&vx=0.1&vy=0.3&m=100&index=1" http://localhost:8081/insertStar
curl -X POST --data "x=200&y=200&vx=0.1&vy=0.3&m=100&index=1" http://localhost:8081/insertStar
curl -X POST --data "x=400&y=400&vx=0.1&vy=0.3&m=100&index=1" http://localhost:8081/insertStar
curl -X GET http://localhost:8081/starlist/go
curl -X GET http://localhost:8081/starlist/csv
