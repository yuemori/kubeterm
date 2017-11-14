NAME = kubeterm

bin/$(NAME):
	go build -o $@ main.go

run:
	go run main.go
