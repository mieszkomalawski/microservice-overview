# Docker Compose - Quick Start

Ten plik zawiera konfigurację Docker Compose do szybkiego uruchomienia aplikacji lokalnie.

## Wymagania

- Docker
- Docker Compose

## Uruchomienie

```bash
# Zbuduj i uruchom wszystkie serwisy
docker-compose up -d

# Zbuduj i uruchom (z logami)
docker-compose up --build

# Zatrzymaj serwisy
docker-compose down

# Zatrzymaj i usuń wolumeny (usuwa dane z bazy)
docker-compose down -v
```

## Dostęp do aplikacji

Po uruchomieniu aplikacja będzie dostępna pod adresem:
- **Frontend**: http://localhost:8080
- **API**: http://localhost:8080/api/

## Dostęp do bazy danych

PostgreSQL jest dostępny na porcie 5432:
- Host: localhost
- Port: 5432
- User: postgres
- Password: postgres
- Database: microservice_overview

Możesz połączyć się używając:
```bash
docker exec -it microservice_overview_db psql -U postgres -d microservice_overview
```

## Sprawdzanie logów

```bash
# Logi aplikacji
docker-compose logs app

# Logi bazy danych
docker-compose logs postgres

# Wszystkie logi
docker-compose logs -f
```

## Rebuild po zmianach w kodzie

```bash
docker-compose up --build -d
```

