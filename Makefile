.PHONY: build run.create

GIT_SHORT_COMMIT=$(shell git log -1 --pretty=format:%h)

build:
	go build -o dist -ldflags='-X "github.com/thibaultmg/clingua/cmd/clingua/cmd.ClinguaVersion=v0.0.0-$(GIT_SHORT_COMMIT)"' ./cmd/clingua

run.create: build
	./dist/clingua -c $(PWD)/resources/.clingua.yaml create ace
