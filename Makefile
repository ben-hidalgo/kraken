docker-mysql:
	docker-compose up --build mysql

docker-grant:
	docker exec -it mysql mysql -prootpass -e "GRANT ALL ON *.* TO 'local'@'%';"

run-user:
	GO111MODULE=on go run services/user/main.go

run-scene:
	GO111MODULE=on go run services/scene/main.go

run-rproxy:
	GO111MODULE=on go run services/rproxy/main.go
