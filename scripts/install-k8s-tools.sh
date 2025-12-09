#!/bin/bash

set -e

echo "🔧 Instalacja narzędzi Kubernetes do lintowania"
echo "================================================"
echo ""

# Sprawdź czy kubectl jest zainstalowany
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl nie jest zainstalowany"
    echo "Zainstaluj kubectl: https://kubernetes.io/docs/tasks/tools/"
    exit 1
else
    echo "✅ kubectl jest zainstalowany: $(kubectl version --client --short)"
fi

# Instalacja kind
if ! command -v kind &> /dev/null; then
    echo ""
    echo "📦 Instalowanie kind..."
    curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
    chmod +x ./kind
    
    # Spróbuj zainstalować systemowo, jeśli nie to lokalnie
    if sudo -n true 2>/dev/null; then
        sudo mv ./kind /usr/local/bin/kind
        echo "✅ kind zainstalowany systemowo"
    else
        mkdir -p ~/.local/bin
        mv ./kind ~/.local/bin/kind
        export PATH="$HOME/.local/bin:$PATH"
        echo "✅ kind zainstalowany lokalnie w ~/.local/bin"
        echo "   Dodaj do PATH: export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi
else
    echo "✅ kind jest zainstalowany: $(kind --version)"
fi

# Instalacja kube-score
if ! command -v kube-score &> /dev/null; then
    echo ""
    echo "📦 Instalowanie kube-score..."
    wget -q https://github.com/zegl/kube-score/releases/download/v1.18.0/kube-score_1.18.0_linux_amd64.tar.gz
    tar xf kube-score_1.18.0_linux_amd64.tar.gz
    rm -f kube-score_1.18.0_linux_amd64.tar.gz
    
    # Spróbuj zainstalować systemowo, jeśli nie to lokalnie
    if sudo -n true 2>/dev/null; then
        sudo mv kube-score /usr/local/bin/kube-score
        echo "✅ kube-score zainstalowany systemowo"
    else
        mkdir -p ~/.local/bin
        mv kube-score ~/.local/bin/kube-score
        export PATH="$HOME/.local/bin:$PATH"
        echo "✅ kube-score zainstalowany lokalnie w ~/.local/bin"
        echo "   Dodaj do PATH: export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi
else
    echo "✅ kube-score jest zainstalowany: $(kube-score version)"
fi

echo ""
echo "✅ Wszystkie narzędzia są gotowe!"
echo ""
echo "Możesz teraz uruchomić: ./scripts/lint-k8s.sh"

