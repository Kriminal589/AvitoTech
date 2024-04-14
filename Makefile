build:
	docker-compose -f docker-compose.yml up -d --build

integration_test:
	docker compose --file docker-compose.test.yml up -d --build
	go test -v './integration_tests/get_user_banner/get_user_banner_test.go'