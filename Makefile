OUT_BINARY=imagesApp

build:
	env GOOS=linux CGO_ENABLED=0 go build -o ${OUT_BINARY} ./
 