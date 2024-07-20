all:
	mkdir -p ./bin
	go build -o emrs ./cmd/cli/*.go && mv emrs ./bin/emrs

run:
	go run ./cmd/cli/*.go

clean:
	rm -fvx ./bin/emrs

