VERSION ?= v$(shell cat .version)

up: tidy
	go run ./__example__/

tidy:
	go mod tidy

version:
	git add . && git commit -m "chore: bump version to $(VERSION)" && git tag $(VERSION) && git push origin $(VERSION)