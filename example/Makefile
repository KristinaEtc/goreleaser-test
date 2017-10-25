all:  # run
	echo "MAKEFILE"
	
run: 
	go build main.go


.PHONY: dep
prog: dep
	go build main.go

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
