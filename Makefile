#defaults
port ?= 25
name ?= smtp_listener
env = ENV_POD
network ?= my_network

SMTP_ADDRESS  ?= smtp_listener:25
SMTP_DOMAIN   ?= smtp_listener
AMQP_URL	  ?= amqp://guest:guest@some-rabbit:5672
AMQP_QUEUE	  ?= smtp-queue
AMQP_EXCHANGE ?= smtp-exchange

#для докера
d_connect:
#как обычный флаг d_connect container=1234
	docker exec -it ${container} /bin/bash

d_run:
	docker run --rm -p ${port}:25 -e ENV=$(env) \
	--name ${name} \
	-e SMTP_ADDRESS=${SMTP_ADDRESS} \
	-e SMTP_DOMAIN=${SMTP_DOMAIN} \
	-e AMQP_URL=${AMQP_URL} \
	-e AMQP_QUEUE=${AMQP_QUEUE} \
	-e AMQP_EXCHANGE=${AMQP_EXCHANGE} \
	--network ${network} \
	end1essrage/listener-smtp-rmq:latest

d_build: 
	docker build -t end1essrage/listener-smtp-rmq .

d_push:
	docker push end1essrage/listener-smtp-rmq:${tag}

#для подмена
p_connect:
#как обычный флаг p_connect container=1234
	podman exec -it ${container} /bin/bash
	
p_run:
	podman run -p ${port}:25 -e ENV=$(env) \
	-e SMTP_ADDRESS=${SMTP_ADDRESS} \
	-e SMTP_DOMAIN=${SMTP_DOMAIN} \
	-e AMQP_URL=${AMQP_URL} \
	-e AMQP_QUEUE=${AMQP_QUEUE} \
	-e AMQP_EXCHANGE=${AMQP_EXCHANGE} \
	--network ${network} \
	end1essrage/listener-smtp-rmq:latest

p_build: 
	podman build -t end1essrage/listener-smtp-rmq .

p_push:
	podman push end1essrage/listener-smtp-rmq:${tag}