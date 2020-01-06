.PHONY:	serv web

prod:
	docker build -t ikaven123/redis-admin .
serv:
	cd serv && go build -o redis-admin
web:
	docker run --rm -v `pwd`/web:/building -w /building \
		-e NPM_CONFIG_REGISTRY=https://registry.npm.taobao.org \
		node:8 bash -c "npm install && npm run build:prod"
build:
	docker build -f Dockerfile.stage . -t redis-admin:stage
run:
	-docker rm -f redis-admin
	docker run --name redis-admin --rm -d -p 6789:6789 \
		-v /data/redis-admin:/usr/share/redis-admin/data \
		redis-admin:stage
