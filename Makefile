run:
	go run cmd/video-to-ascii/main.go -file ./sample-data/IMG_1135.MOV

build:
	go build -o bin/video-to-ascii cmd/video-to-ascii/main.go
