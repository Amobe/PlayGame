
verify:
	go vet ./... && staticcheck ./... && go test ./...

docker.build:
	docker buildx build --platform linux/arm64/v8 -t playgame/server:latest .

docker.debug.build:
	docker buildx build --platform linux/arm64/v8 -t playgame/server:debug -f Dockerfile.debug .

docker.cloud.build:
	docker buildx build --platform linux/amd64 -t playgame/server:latest-k8s .

docker.cloud.push:
	docker tag playgame/server:latest-k8s asia-east1-docker.pkg.dev/playgame-373510/playgame/server:latest-k8s
	docker push asia-east1-docker.pkg.dev/playgame-373510/playgame/server:latest-k8s
