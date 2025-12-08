#!/bin/bash

set -e

echo "🚀 Wdrażanie Microservice Overview do Kubernetes..."

# Sprawdź czy kubectl jest dostępny
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl nie jest zainstalowany"
    exit 1
fi

# Sprawdź czy jesteśmy połączeni z clusterem
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Nie można połączyć się z Kubernetes clusterem"
    exit 1
fi

echo "📦 Tworzenie namespace..."
kubectl apply -f namespace.yaml

echo "🗄️  Wdrażanie PostgreSQL..."
kubectl apply -f postgres-secret.yaml
kubectl apply -f postgres-configmap.yaml
kubectl apply -f postgres-pvc.yaml
kubectl apply -f postgres-deployment.yaml
kubectl apply -f postgres-service.yaml

echo "⏳ Oczekiwanie na gotowość PostgreSQL..."
kubectl wait --for=condition=ready pod -l app=postgres -n microservice-overview --timeout=120s

echo "📱 Wdrażanie aplikacji..."
kubectl apply -f app-configmap.yaml
kubectl apply -f app-secret.yaml
kubectl apply -f app-deployment.yaml
kubectl apply -f app-service.yaml

echo "🌐 Wdrażanie Ingress..."
kubectl apply -f ingress.yaml

echo "⏳ Oczekiwanie na gotowość aplikacji..."
kubectl wait --for=condition=ready pod -l app=microservice-overview -n microservice-overview --timeout=120s

echo ""
echo "✅ Wdrożenie zakończone!"
echo ""
echo "📊 Status zasobów:"
kubectl get all -n microservice-overview

echo ""
echo "🌐 Aby uzyskać dostęp do aplikacji:"
echo ""
echo "1. Dodaj do /etc/hosts (wymaga sudo):"
echo "   echo '127.0.0.1 microservice-overview.local' | sudo tee -a /etc/hosts"
echo ""
echo "2. Lub użyj port-forward:"
echo "   kubectl port-forward -n microservice-overview service/microservice-overview 8080:80"
echo ""
echo "3. Następnie otwórz w przeglądarce:"
echo "   http://microservice-overview.local (z Ingress)"
echo "   lub http://localhost:8080 (z port-forward)"
echo ""

