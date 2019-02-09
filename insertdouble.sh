#!/usr/bin/env bash
curl -X POST http://localhost:8081/deleteStars
curl -X POST http://localhost:8081/deleteNodes 
curl -X POST --data "w=1000" http://localhost:8081/new
curl -X POST --data "x=750&y=250&vx=0.1&vy=0.3&m=100&index=1" http://localhost:8081/insertStar
curl -X POST --data "x=250&y=250&vx=0.1&vy=0.3&m=100&index=1" http://localhost:8081/insertStar
