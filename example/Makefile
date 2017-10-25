all:  # run
	echo "====================FROM MAKEFILE"
	
run: 
	go build main.go
	#./experimental

.PHONY: dep
prog: dep
	echo "====================MAKEFILE"
	go build main.go
	#./main

dep:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: docker

docker:
	docker build -t with_mold .
	docker run -d with_mold

clean:
	rm -rf $(NAME) vendor/* dist coverage.out
	go clean -i ./...

l: 
	ls
