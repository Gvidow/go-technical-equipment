build:
	go build -o bin/app cmd/app/*.go

run: build
	./bin/app

addbin:
	export PATH=$PATH:/home/artem/go/bin

swaggen:
	swag init -g cmd/app/main.go

update_swag:
	swag init -g cmd/app/main.go
	go build -o bin/app cmd/app/*.go
	./bin/app
