.PHONY: clean

build/scf: cmd/tc_scf/main.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $@ $<

build/scf.zip: build/scf
	zip -j $@ $<

clean:
	rm build/*
