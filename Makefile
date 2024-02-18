DOCKER_BUILD=docker build
DOCKER_RUN=docker run -d
DOCKER_NETWORK_CREATE=docker network create
DOCKER_RM=docker rm -f
DOCKER_NETWORK_RM=docker network rm
NETWORK_NAME=wordofwisdom_network
SERVER_IMAGE=wordofwisdom-server
CLIENT_IMAGE=wordofwisdom-client
SERVER_CONTAINER_NAME=wowsrv
CLIENT_CONTAINER_NAME=wowclient

all: clean build_network build_server run_server build_client run_client

build_network:
	-$(DOCKER_NETWORK_CREATE) --driver bridge $(NETWORK_NAME)

build_server:
	cd server && $(DOCKER_BUILD) -t $(SERVER_IMAGE) .

build_client:
	cd client && $(DOCKER_BUILD) -t $(CLIENT_IMAGE) .

run_server:
	docker run -d --name $(SERVER_CONTAINER_NAME) --network $(NETWORK_NAME) $(SERVER_IMAGE)

run_client:
	docker run --rm --name $(CLIENT_CONTAINER_NAME) --network $(NETWORK_NAME) $(CLIENT_IMAGE)

clean:
	-$(DOCKER_RM) $(SERVER_CONTAINER_NAME) || true
	-$(DOCKER_RM) $(CLIENT_CONTAINER_NAME) || true
	-$(DOCKER_NETWORK_RM) $(NETWORK_NAME) || true

.PHONY: all build_network build_server build_client run_server run_client clean