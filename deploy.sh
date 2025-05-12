#!/bin/bash
set -e

curl -L https://istio.io/downloadIstio | sh -
cd istio-*
export PATH=$PWD/bin:$PATH
istioctl install -y
kubectl label namespace default istio-injection=enabled
cd ..

eval $(minikube docker-env)
docker build -t server:v1.0 .
docker pull alpine/curl

kubectl apply -f gateway.yaml -f virtualservice.yaml -f destinationrule.yaml

kubectl apply -f configmap.yaml -f deployment.yaml -f service.yaml -f daemonset.yaml -f cronjob.yaml

echo "Ожидание Deployment"
kubectl rollout status deployment/server-deployment

echo "Ожидание DaemonSet"
kubectl rollout status daemonset/server-daemonset

echo "Установка Prometheus"
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
kubectl create namespace monitoring || true
helm install kube-prom-stack prometheus-community/kube-prometheus-stack -n monitoring

echo "Ожидание Prometheus"
sleep 900

kubectl apply -f servicemonitor.yaml
kubectl apply -f requestavgrule.yaml
kubectl port-forward -n istio-system svc/istio-ingressgateway 8080:80
kubectl port-forward -n monitoring svc/kube-prom-stack-kube-prome-prometheus 9090:9090