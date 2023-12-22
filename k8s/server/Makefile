.PHONY: server
server:
	-kubectl delete deployment gint-server-deployment
	-kubectl delete service gint-server-service
	kubectl apply -f k8s-gint-deployment.yaml
	kubectl apply -f k8s-gint-service.yaml

.PHONY: server-del
server-del:
	-kubectl delete deployment gint-server-deployment
	-kubectl delete service gint-server-service