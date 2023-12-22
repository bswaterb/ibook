.PHONY: mysql
mysql:
	-kubectl delete deployment gint-mysql-deployment
	-kubectl delete service gint-mysql-service
	-kubectl patch pv mysql-volume-1 -p '{"metadata":{"finalizers":null}}'
	# -kubectl delete pv mysql-volume-1
	-kubectl patch pv mysql-volume-1gb -p '{"metadata":{"finalizers":null}}'
	# -kubectl delete pvc mysql-volume-1gb
	kubectl apply -f k8s-mysql-pv.yaml
	kubectl apply -f k8s-mysql-pvc.yaml
	kubectl apply -f k8s-mysql-deployment.yaml
	kubectl apply -f k8s-mysql-service.yaml

.PHONY: mysql-del
mysql-del:
	-kubectl delete deployment gint-mysql-deployment
	-kubectl delete service gint-mysql-service
	-kubectl patch pv mysql-volume-1 -p '{"metadata":{"finalizers":null}}'
	-kubectl patch pv mysql-volume-1gb -p '{"metadata":{"finalizers":null}}'
	-kubectl delete pvc mysql-volume-1gb
	-kubectl delete pv mysql-volume-1


