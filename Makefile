.PHONY: build
build:
	go build -o build/planturl main.go

.PHONY: clean
clean:
	rm -rf build/

