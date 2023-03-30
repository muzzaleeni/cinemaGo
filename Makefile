EXE=my-project

build:
	go build -o $(EXE) ./cmd/my-project/main.go 

run: build
	./$(EXE)
