
docker.build:
	docker buildx build --platform linux/arm64/v8 -t playgame/server:latest .

docker.push:
	docker tag playgame/server:latest asia-east1-docker.pkg.dev/playgame-373510/playgame/server:latest
	docker push asia-east1-docker.pkg.dev/playgame-373510/playgame/server:latest
