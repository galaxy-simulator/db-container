curl -X POST http://localhost:8081/deleteStars 
curl -X POST http://localhost:8081/deleteNodes 
curl -X POST --data "w=1000" http://localhost:8081/new
curl -X POST --data "filename=200000.csv" http://localhost:8081/insertList
