doc:
	swag init -g cmd/server/main.go -s ./../trash/docs
docker:
	docker build -t cloud-disk:latest -f Dockerfile.server .
	docker build -t cloud-disk_worker:latest -f Dockerfile.worker .
docker_upload: docker
	docker tag cloud-disk:latest registry.cn-hangzhou.aliyuncs.com/wqer1019/cloud-disk:latest
	docker push registry.cn-hangzhou.aliyuncs.com/wqer1019/cloud-disk:latest
	docker tag cloud-disk_worker:latest registry.cn-hangzhou.aliyuncs.com/wqer1019/cloud-disk_worker:latest
	docker push registry.cn-hangzhou.aliyuncs.com/wqer1019/cloud-disk_worker:latest
run:
	docker run -d cloud-disk:latest
	docker run -d cloud-disk_worker:latest