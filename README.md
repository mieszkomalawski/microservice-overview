# Microservice Overview - Wizualizacja zależności między mikroserwisami

Aplikacja do wizualizacji zależności między mikroserwisami w formie grafu.

## Funkcjonalności

- **RESTful API HTTP** do zarządzania wierzchołkami (mikroserwisami) i relacjami między nimi
- **CRUD** dla wierzchołków i relacji
- **Wizualizacja grafu** w przeglądarce
- **Storage**: PostgreSQL (produkcja) lub tryb developerski w pamięci

## Instalacja

```bash
go mod download
```

## Testy

### Uruchomienie testów lokalnie

```bash
# Wszystkie testy
go test ./...

# Tylko unit testy
go test ./storage/...

# Tylko testy integracyjne
go test ./handlers/...

# Z pokryciem kodu
go test -cover ./...
```

### CI/CD

Projekt używa GitHub Actions do automatycznego uruchamiania testów i lintowania:
- **Testy Go**: uruchamiane automatycznie po każdym push do branchy `main`, `master` lub `develop`
- **Testy są uruchamiane również dla Pull Requestów**
- **Testy są uruchamiane dla wersji Go 1.23 i 1.24**
- **Raport pokrycia kodu** jest zapisywany jako artifact
- **Lintowanie Kubernetes**: automatyczne sprawdzanie konfiguracji K8s:
  - Walidacja schematów z `kubeval`
  - Walidacja z `kubectl --dry-run`
  - Sprawdzanie best practices i bezpieczeństwa z `kube-score`

## Konfiguracja

Aplikacja używa zmiennych środowiskowych:

- `DB_HOST` - host bazy danych PostgreSQL (domyślnie: localhost)
- `DB_PORT` - port bazy danych (domyślnie: 5432)
- `DB_USER` - użytkownik bazy danych (domyślnie: postgres)
- `DB_PASSWORD` - hasło bazy danych (domyślnie: postgres)
- `DB_NAME` - nazwa bazy danych (domyślnie: microservice_overview)
- `DEV_MODE` - tryb developerski w pamięci (domyślnie: false)

## Uruchomienie

### Tryb produkcyjny (PostgreSQL)

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=microservice_overview
go run main.go
```

### Tryb developerski (w pamięci)

```bash
export DEV_MODE=true
go run main.go
```

Aplikacja będzie dostępna pod adresem: http://localhost:8080

- **API**: http://localhost:8080/api/
- **Frontend**: http://localhost:8080/

## Uruchomienie z Docker Compose

Najprostszy sposób na uruchomienie aplikacji z wszystkimi zależnościami:

```bash
# Zbuduj i uruchom
docker-compose up -d

# Zobacz logi
docker-compose logs -f

# Zatrzymaj
docker-compose down
```

Aplikacja będzie dostępna pod adresem: http://localhost:8080

Szczegóły w pliku [docker-compose.README.md](docker-compose.README.md)

## Uruchomienie w Kubernetes

Aby uruchomić aplikację w Kubernetes z lokalną domeną:

### Wymagania
- Kubernetes cluster (minikube, kind, k3d)
- kubectl
- Nginx Ingress Controller 

```bash
minikube addons enable ingress
```
 
### Szybkie wdrożenie

```bash
# 1. Zbuduj obraz Docker (dla minikube)
eval $(minikube docker-env)
docker build -t microservice-overview:latest .
eval $(minikube docker-env -u)

# 2. Wdróż wszystko
cd k8s
./deploy.sh

# 3. Skonfiguruj lokalną domenę
echo "127.0.0.1 microservice-overview.local" | sudo tee -a /etc/hosts
```

Dodatkowo przydaje się minikube dashboard

```bash
minikube dashboard
```


Aplikacja będzie dostępna pod adresem: http://microservice-overview.local

Szczegóły w katalogu [k8s/README.md](k8s/README.md)

## API Endpoints

### Wierzchołki (Mikroserwisy)
- `GET /api/vertices` - Lista wszystkich wierzchołków
- `GET /api/vertices/:id` - Pobierz wierzchołek po ID
- `POST /api/vertices` - Utwórz nowy wierzchołek
- `PUT /api/vertices/:id` - Aktualizuj wierzchołek
- `DELETE /api/vertices/:id` - Usuń wierzchołek

### Relacje (Połączenia)
- `GET /api/edges` - Lista wszystkich relacji
- `GET /api/edges/:id` - Pobierz relację po ID
- `POST /api/edges` - Utwórz nową relację
- `PUT /api/edges/:id` - Aktualizuj relację
- `DELETE /api/edges/:id` - Usuń relację

### Graf
- `GET /api/graph` - Pobierz pełny graf (wszystkie wierzchołki i relacje)

## Kolekcja Postman

Gotowa kolekcja Postman z wszystkimi endpointami i przykładami jest dostępna w pliku `postman_collection.json`.

### Import do Postman

1. Otwórz Postman
2. Kliknij **Import** (lub `Ctrl+O` / `Cmd+O`)
3. Wybierz plik `postman_collection.json`
4. Kolekcja zostanie zaimportowana z wszystkimi requestami

### Konfiguracja

Kolekcja używa zmiennej `base_url` z domyślną wartością `http://localhost:8080`. Możesz ją zmienić w:
- Postman → Variables → `base_url`

Dla Kubernetes z lokalną domeną ustaw:
- `http://microservice-overview.local`

### Zawartość kolekcji

- **Vertices (Wierzchołki)**: wszystkie operacje CRUD + przykłady tworzenia różnych serwisów
- **Edges (Relacje)**: wszystkie operacje CRUD + przykłady różnych typów relacji
- **Graph (Graf)**: pobieranie pełnego grafu

## Format danych

### Wierzchołek (Vertex)
```json
{
  "id": "string",
  "name": "string",
  "description": "string (opcjonalne)",
  "parent_id": "string (opcjonalne, ID rodzica dla hierarchii)"
}
```

**Hierarchia wierzchołków:**
- Wierzchołki mogą zawierać w sobie inne wierzchołki poprzez pole `parent_id`
- Wierzchołek bez `parent_id` (lub `null`) jest na najwyższym poziomie (root)
- Wierzchołek z `parent_id` jest potomkiem innego wierzchołka
- **Ważne**: Wierzchołek nie może mieć połączeń (edges) z wierzchołkami w sobie
- **Połączenia mogą istnieć tylko między wierzchołkami najniższego poziomu** (bez dzieci)

### Relacja (Edge)
```json
{
  "id": "string",
  "from": "string (ID wierzchołka źródłowego)",
  "to": "string (ID wierzchołka docelowego)",
  "type": "string (opcjonalne, typ relacji)"
}
```
