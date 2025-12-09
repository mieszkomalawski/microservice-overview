#!/bin/bash

set -e

echo "🔍 Kubernetes Linter - Local Execution"
echo "======================================"
echo ""

# Sprawdź czy kubectl jest zainstalowany
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl nie jest zainstalowany"
    echo "Zainstaluj kubectl: https://kubernetes.io/docs/tasks/tools/"
    exit 1
fi

# Sprawdź czy kind jest zainstalowany
if ! command -v kind &> /dev/null; then
    echo "❌ kind nie jest zainstalowany"
    echo "Zainstaluj kind: https://kind.sigs.k8s.io/docs/user/quick-start/#installation"
    exit 1
fi

# Sprawdź czy kube-score jest zainstalowany
KUBE_SCORE_BIN="./kube-score"
if ! command -v kube-score &> /dev/null; then
    if [ ! -f "$KUBE_SCORE_BIN" ]; then
        echo "⚠️  kube-score nie jest zainstalowany, instaluję lokalnie..."
        wget -q https://github.com/zegl/kube-score/releases/download/v1.18.0/kube-score_1.18.0_linux_amd64.tar.gz
        tar xf kube-score_1.18.0_linux_amd64.tar.gz
        rm -f kube-score_1.18.0_linux_amd64.tar.gz
        echo "✅ kube-score zainstalowany lokalnie"
    fi
    KUBE_SCORE_BIN="./kube-score"
else
    KUBE_SCORE_BIN="kube-score"
fi

echo ""
echo "🚀 Tworzenie lokalnego klastra Kubernetes (kind)..."
echo ""

# Sprawdź czy klaster już istnieje
if kind get clusters | grep -q "test-cluster"; then
    echo "ℹ️  Klaster test-cluster już istnieje, używam istniejącego"
    kind get kubeconfig --name test-cluster > /tmp/kind-kubeconfig
    export KUBECONFIG=/tmp/kind-kubeconfig
else
    echo "Tworzenie nowego klastra..."
    kind create cluster --name test-cluster --wait 5m
    kind get kubeconfig --name test-cluster > /tmp/kind-kubeconfig
    export KUBECONFIG=/tmp/kind-kubeconfig
fi

echo ""
echo "📋 Walidacja manifestów Kubernetes..."
echo ""

# Utwórz namespace jeśli istnieje
if [ -f "k8s/namespace.yaml" ]; then
    echo "Utwarzanie namespace z k8s/namespace.yaml"
    kubectl apply -f k8s/namespace.yaml
fi

# Waliduj wszystkie pliki
VALIDATION_FAILED=0
for file in k8s/*.yaml; do
    if [ -f "$file" ] && [ "$(basename "$file")" != "namespace.yaml" ]; then
        echo "Walidacja: $file"
        if ! kubectl apply --dry-run=server -f "$file" 2>&1; then
            echo "❌ Błąd walidacji: $file"
            VALIDATION_FAILED=1
        else
            echo "✅ $file - OK"
        fi
        echo ""
    fi
done

echo ""
echo "🔍 Sprawdzanie best practices i bezpieczeństwa (kube-score)..."
echo ""

$KUBE_SCORE_BIN score k8s/*.yaml || true

echo ""
echo "🔒 Sprawdzanie krytycznych problemów bezpieczeństwa..."
echo ""

$KUBE_SCORE_BIN score --output-format ci k8s/*.yaml > /tmp/kube-score-report.txt || true

# Filtruj akceptowalne problemy dla postgres (niski UID/GID i writable filesystem są wymagane)
if grep -q "CRITICAL\|WARNING" /tmp/kube-score-report.txt; then
    # Sprawdź czy są tylko akceptowalne problemy z postgres
    CRITICAL_COUNT=$(grep -c "CRITICAL" /tmp/kube-score-report.txt || echo "0")
    POSTGRES_ACCEPTABLE=$(grep -c "postgres.*low user ID\|postgres.*low group ID\|postgres.*writable root filesystem" /tmp/kube-score-report.txt || echo "0")
    
    if [ "$CRITICAL_COUNT" -eq "$POSTGRES_ACCEPTABLE" ]; then
        echo "✅ Brak krytycznych problemów bezpieczeństwa (problemy z postgres są akceptowalne - wymagane przez bazę danych)"
    else
        echo "⚠️  Znaleziono problemy bezpieczeństwa lub ostrzeżenia:"
        cat /tmp/kube-score-report.txt
        VALIDATION_FAILED=1
    fi
else
    echo "✅ Brak krytycznych problemów bezpieczeństwa"
fi

echo ""
echo "🧹 Czyszczenie..."

# Opcjonalnie usuń klaster (odkomentuj jeśli chcesz)
# kind delete cluster --name test-cluster

if [ $VALIDATION_FAILED -eq 1 ]; then
    echo ""
    echo "❌ Walidacja zakończona z błędami"
    exit 1
else
    echo ""
    echo "✅ Wszystkie walidacje przeszły pomyślnie!"
    exit 0
fi

