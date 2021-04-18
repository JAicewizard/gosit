ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

clean:
	rm positrs/positrs.h

positrs/positrs.h: $(wildcard positrs/**/*.rs) $(wildcard rpositrss/**/*.toml)
	cd positrs; ~/.cargo/bin/cbindgen --config cbindgen.toml --crate positrs --output positrs.h --lang c

positrs/libpositrs: $(wildcard positrs/**/*.rs) $(wildcard rpositrss/**/*.toml)
	cd positrs; cargo build --release

rustfuzzing: *.go positrs/positrs.h positrs/libpositrs
	cp positrs/target/release/libpositrs.so ./positrs
	go test -tags cgo_bench -ldflags="-r $(ROOT_DIR)positrs" --run=Fuz -parallel 4 -timeout 30m
