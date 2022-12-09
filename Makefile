runES:
	docker run -e ES_JAVA_OPTS="-Xms1g -Xmx1g" -e "discovery.type=single-node" -e "xpack.security.enabled=false" --net elastic -p 9200:9200 docker.elastic.co/elasticsearch/elasticsearch:8.5.2

run:
	go run main.go