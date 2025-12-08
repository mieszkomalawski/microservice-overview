# Kubernetes Deployment Guide

Ten katalog zawiera konfigurację Kubernetes do uruchomienia aplikacji Microservice Overview lokalnie.

## Wymagania

- Kubernetes cluster (minikube, kind, k3d, lub lokalny cluster)
- kubectl skonfigurowany do komunikacji z clusterem
- Docker do budowania obrazu aplikacji
- Nginx Ingress Controller (dla lokalnej domeny)

## Instalacja Nginx Ingress Controller

### Minikube
```bash
minikube addons enable ingress
```

### Kind/K3d lub inny cluster
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml
```

## Budowanie i ładowanie obrazu Docker

### Dla Minikube
```bash
# Ustaw docker environment na minikube
eval $(minikube docker-env)

# Zbuduj obraz
docker build -t microservice-overview:latest .

# Reset docker environment
eval $(minikube docker-env -u)
```

### Dla Kind
```bash
# Zbuduj obraz
docker build -t microservice-overview:latest .

# Załaduj do kind
kind load docker-image microservice-overview:latest
```

### Dla K3d
```bash
# Zbuduj obraz
docker build -t microservice-overview:latest .

# Załaduj do k3d (zastąp CLUSTER_NAME nazwą swojego clusteru)
kubectl get nodes -o wide  # znajdź nazwę node'a
docker save microservice-overview:latest | k3d image import -c CLUSTER_NAME
```

## Wdrożenie

1. Utwórz namespace:
```bash
kubectl apply -f namespace.yaml
```

2. Wdróż PostgreSQL:
```bash
kubectl apply -f postgres-secret.yaml
kubectl apply -f postgres-configmap.yaml
kubectl apply -f postgres-pvc.yaml
kubectl apply -f postgres-deployment.yaml
kubectl apply -f postgres-service.yaml
```

3. Wdróż aplikację:
```bash
kubectl apply -f app-configmap.yaml
kubectl apply -f app-secret.yaml
kubectl apply -f app-deployment.yaml
kubectl apply -f app-service.yaml
```

4. Wdróż Ingress:
```bash
kubectl apply -f ingress.yaml
```

## Konfiguracja lokalnej domeny

### Opcja 1: Dodaj do /etc/hosts (Linux/Mac)
```bash
# Znajdź IP Ingress Controller
kubectl get ingress -n microservice-overview

# Dodaj do /etc/hosts (wymaga sudo)
echo "127.0.0.1 microservice-overview.local" | sudo tee -a /etc/hosts
```

### Opcja 2: Użyj nip.io (automatyczne DNS)
Edytuj `ingress.yaml` i odkomentuj sekcję z `nip.io`, a następnie użyj adresu:
```
http://microservice-overview.127.0.0.1.nip.io
```

### Opcja 3: Port forwarding (bez Ingress)
```bash
kubectl port-forward -n microservice-overview service/microservice-overview 8080:80
```
Następnie otwórz: http://localhost:8080

## Sprawdzanie statusu

```bash
# Sprawdź wszystkie zasoby
kubectl get all -n microservice-overview

# Sprawdź logi aplikacji
kubectl logs -n microservice-overview -l app=microservice-overview

# Sprawdź logi PostgreSQL
kubectl logs -n microservice-overview -l app=postgres

# Sprawdź status Ingress
kubectl get ingress -n microservice-overview
```

## Usuwanie

```bash
kubectl delete namespace microservice-overview
```

## Alternatywnie - użyj jednego polecenia

Możesz wdrożyć wszystko naraz:
```bash
kubectl apply -f k8s/
```

## Troubleshooting

### Aplikacja nie może połączyć się z bazą danych
```bash
# Sprawdź czy PostgreSQL jest gotowy
kubectl get pods -n microservice-overview

# Sprawdź logi aplikacji
kubectl logs -n microservice-overview -l app=microservice-overview
```

### Ingress nie działa
```bash
# Sprawdź czy Ingress Controller jest uruchomiony
kubectl get pods -n ingress-nginx

# Sprawdź konfigurację Ingress
kubectl describe ingress -n microservice-overview
```

### Problem z PersistentVolume
Jeśli używasz minikube, upewnij się że masz włączony storage provisioner:
```bash
minikube addons enable default-storageclass
```

