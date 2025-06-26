run:
	go run cmd/video-to-ascii/main.go -file ./sample-data/IMG_1033.MOV

debug:
	go run cmd/video-to-ascii/main.go -debug -file ./sample-data/IMG_1033.MOV

build:
	go build -o bin/video-to-ascii cmd/video-to-ascii/main.go
