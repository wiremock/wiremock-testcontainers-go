.PHONY: docker
docker:
	docker build -t wiremock/wiremock:2.35.0-for-tc .
