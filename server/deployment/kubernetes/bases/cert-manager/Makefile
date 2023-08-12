all:

get-cert-manager:
	wget -O cert-manager.origin.yaml https://github.com/cert-manager/cert-manager/releases/download/v1.12.3/cert-manager.yaml
	# Following the instruction to fix the issue.
	# https://cert-manager.io/docs/installation/compatibility/#gke-autopilot
	sed -e 's/kube-system/cert-manager/g' cert-manager.origin.yaml > cert-manager.yaml
	rm -f cert-manager.origin.yaml
