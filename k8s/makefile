SHELL := /bin/bash

export PROJECT = sugar
CLUSTER = sugar
DATABASE = sugar-db
ZONE = europe-west1-b

config:
	@echo Setting environment for $(PROJECT)
	gcloud config set project $(PROJECT)
	gcloud config set compute/zone $(ZONE)
	gcloud auth configure-docker
	@echo ======================================================================

project:
	gcloud projects create $(PROJECT)
	gcloud beta billing projects link $(PROJECT) --billing-account=$(ACCOUNT_ID)
	gcloud services enable container.googleapis.com
	@echo ======================================================================

cluster:
	gcloud container clusters create $(CLUSTER) --enable-ip-alias --num-nodes=2 --machine-type=n1-standard-2
	gcloud compute instances list
	@echo ======================================================================

upload:
	docker push gcr.io/$(PROJECT)/sugar-api-amd64:1.0
	docker push gcr.io/$(PROJECT)/metrics-amd64:1.0
	@echo ======================================================================

database:
	# Create User/Password
	gcloud beta sql instances create $(DATABASE) --database-version=POSTGRES_9_6 --no-backup --tier=db-f1-micro --zone=$(ZONE) --no-assign-ip --network=default
	gcloud sql instances describe $(DATABASE)
	@echo ======================================================================

db-assign-ip:
	gcloud sql instances patch $(DATABASE) --authorized-networks=[$(PUBLIC-IP)/32]
	gcloud sql instances describe $(DATABASE)
	@echo ======================================================================

db-private-ip:
	# IMPORTANT: Make sure you run this command and get the private IP of the DB.
	gcloud sql instances describe $(DATABASE)
	@echo ======================================================================

services:
	# These scripts needs to be edited for the PROJECT and PRIVATE_DB_IP markers before running.
	kubectl create -f deploy-sugar-api.yaml
	kubectl expose -f expose-sugar-api.yaml --type=LoadBalancer
	@echo ======================================================================

status:
	gcloud container clusters list
	kubectl get nodes
	kubectl get pods
	kubectl get services sugar-api
	@echo ======================================================================

shell:
	# kubectl get pods
	kubectl exec -it <POD NAME> --container sugar-api  -- /bin/sh
	# ./admin --db-disable-tls=1 migrate
	# ./admin --db-disable-tls=1 seed
	@echo ======================================================================

delete:
	kubectl delete services sugar-api
	kubectl delete deployment sugar-api
	gcloud container clusters delete $(CLUSTER)
	gcloud projects delete sugar-api
	gcloud container images delete gcr.io/$(PROJECT)/sugar-api-amd64:1.0 --force-delete-tags
	gcloud container images delete gcr.io/$(PROJECT)/metrics-amd64:1.0 --force-delete-tags
	docker image remove gcr.io/sugar-api/sugar-api-amd64:1.0
	docker image remove gcr.io/sugar-api/metrics-amd64:1.0
	@echo ======================================================================

#===============================================================================
# GKE Installation
#
# Install the Google Cloud SDK. This contains the gcloud client needed to perform
# some operatings
# https://cloud.google.com/sdk/
#
# Installing the K8s kubectl client.
# https://kubernetes.io/docs/tasks/tools/install-kubectl/