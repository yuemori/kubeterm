SRCS := $(shell find . -type f -name '*.go')

NAME = kubeterm

bin/$(NAME): $(SRCS)
	go build -o $@ main.go

run:
	go run main.go
