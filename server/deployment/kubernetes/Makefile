all:

apply.prod:
	kustomize build overlays/prod | kubectl apply -f -

apply.ingress:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

apply.tls:
	kubectl create secret tls playgame-tls-secret --cert=overlays/prod/ca.crt --key=overlays/prod/ca.key

apply.cert-manager:
	#kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.12.3/cert-manager.crds.yaml
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.12.3/cert-manager.yaml

create.ip:
	gcloud compute addresses create playgame-staging-prod --global

grant.permission:
	kubectl create clusterrolebinding cluster-admin-binding \
        --clusterrole=cluster-admin \
        --user=$(gcloud config get-value core/account)

