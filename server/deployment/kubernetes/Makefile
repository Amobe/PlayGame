all:

apply.prod:
	kustomize build overlays/prod | kubectl apply -f -

apply.ingress:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

apply.tls:
	kubectl create secret tls playgame-tls-secret --cert=overlays/prod/ca.crt --key=overlays/prod/ca.key
