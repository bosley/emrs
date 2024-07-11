all:
	mkdir -p ./bin
	go build -o emrs ./cmd/cli/main.go && mv emrs ./bin/emrs

clean:
	rm -fvx ./bin/emrs

