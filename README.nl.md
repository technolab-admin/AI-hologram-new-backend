# AI-HOLOGRAM-NEW-BACKEND
[![en](https://img.shields.io/badge/lang-EN-red.svg)](README.md)
[![nl](https://img.shields.io/badge/lang-NL-blue.svg)](README.nl.md)

*Auteurs: Twan Wolthaus & Aiden van Wijnbergen*

Software voor het uitvoeren van achtergrondprocessen van het AI-Hologram.

Toepassingen:
- Beheren van WebSocket-communicatie
- Proxyen van verzoeken naar externe AI-services
- Afhandelen van 3D-modelgeneratie en downloads

---

## Installatie

Begin met het clonen van de repository clonen:

```sh
git clone https://github.com/technolab-admin/AI-hologram-new-backend.git
cd AI-hologram-new-backend
```
Na het clonen van de repository, bekijk het bestand `.env.example` en maak hiervan een kopie met de naam `.env`.
Plaats dit bestand in dezelfde map en pas de variabelen aan waar nodig.

### Docker

Om de applicatie te draaien wordt het aangeraden Docker te gebruiken.
Voor het installeren van Docker, klik [hier](https://www.docker.com/products/docker-desktop/).

Na installatie van Docker kan de applicatie starten met:
```sh
docker compose up --build
```

Dit start zowel de server als de WebSocket-server op de poorten die zijn ingesteld in het .env-bestand.

### Tests

Om unit-tests voor de WebSocket-logica uit te voeren:
```sh
cd AI-hologram-new-backend
go test ./test/unit_test/
```

## Projectstructuur

- `cmd/server` – Startpunt van de applicatie
- `internal/http` – HTTP-handlers en routing
- `internal/websockets` – WebSocket-server en validatie
- `internal/meshy` – Beheer van interactie met de Meshy API
- `internal/config` – Omgevingsconfiguratie
