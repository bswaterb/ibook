.PHONY: redis
redis:
	-kubectl delete deployment gint-redis-deployment
	-kubectl delete service gint-redis-service
	kubectl apply -f k8s-redis-deployment.yaml
	kubectl apply -f k8s-redis-service.yaml

.PHONY: redis-del
redis-del:
	-kubectl delete deployment gint-redis-deployment
	-kubectl delete service gint-redis-service