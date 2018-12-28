docker:
	docker build -t xxx:latest -f Dockerfile.server .
	docker build -t xxx_worker:latest -f Dockerfile.worker .
docker_upload: docker
	docker tag xxx:latest registry.cn-hangzhou.aliyuncs.com/zm-dev/xxx:latest
	docker push registry.cn-hangzhou.aliyuncs.com/zm-dev/xxx:latest
	docker tag xxx_worker:latest registry.cn-hangzhou.aliyuncs.com/zm-dev/xxx_worker:latest
	docker push registry.cn-hangzhou.aliyuncs.com/zm-dev/xxx_worker:latest
run:
	docker run -d xxx:latest
	docker run -d xxx_worker:latest