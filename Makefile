FORCE:

DOCKER_COMPOSE_DB=docker-compose.yml


up-db:
	docker-compose -f ${DOCKER_COMPOSE_DB} up  -d

down-db:
	docker-compose -f ${DOCKER_COMPOSE_DB} down 

logs-db:
	docker-compose -f ${DOCKER_COMPOSE_DB} logs 