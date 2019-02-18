IMAGE_NAME := "krast76/transmission-exporter"
IMAGE_TAG := "latest"

all:: build

build::
	go get -d
	go build -o transmission-export

docker-build::
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

clean::
	rm transmission-export
