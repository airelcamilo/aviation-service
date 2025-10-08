# Aviation Service

An **Aviation Data Platform** built with **Go**, **PostgreSQL**, and **Redis**, designed to collect airport information from [AviationAPI](https://www.aviationapi.com/) and provide offline access when the external API is unavailable.  
It also fetches **current weather data** from [WeatherAPI](https://www.weatherapi.com/).

---

## 🚀 Features

- Collect and cache **airport data** from AviationAPI  
- Store airport information in **PostgreSQL** for offline use  
- Retrieve **real-time weather data** for airports  
- CRUD operations for airport records  
- Synchronize incomplete airport data (`status = "PENDING"`)  
- API caching with **Redis**  
- Automated **background sync scheduler**

---

## 🧩 Tech Stack

| Component      | Technology        |
|----------------|------------------|
| Language       | Go (Golang)      |
| Database       | PostgreSQL       |
| Cache          | Redis            |
| External APIs  | AviationAPI, WeatherAPI |
| Deployment     | Docker Compose   |

---

## 🧰 Prerequisites

Make sure you have installed:
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- (optional) [Go](https://golang.org/) if you want to run locally without Docker

---

## 💻 How to Run

### 1. Run the app using Docker

```bash
docker compose up --build
```

### 2️. Run database migrations
```bash
./migrate
```

### 3️. Run the scheduler (for syncing pending airports)
```bash
./schedule
```

---

## 🌐 API Endpoints

### ✈️ Airport Service

| Method     | Endpoint                                                               | Description                                                     |
| ---------- | ---------------------------------------------------------------------- | --------------------------------------------------------------- |
| **GET**    | `/airport?page=1&pageSize=10`                                          | Get all airports with pagination                                |
| **GET**    | `/airport/{id}`                                                        | Get airport by ID                                               |
| **GET**    | `/airport/search?icao=KADT&facilityName=washington&page=1&pageSize=10` | Search airports by ICAO or facility name                        |
| **POST**   | `/airport`                                                             | Create new airport record. If incomplete, status = `"PENDING"`. |
| **PUT**    | `/airport/{id}`                                                        | Update airport by ID                                            |
| **DELETE** | `/airport/{id}`                                                        | Delete airport by ID                                            |

Example `POST` Body

```json
{
    "icao_ident": "KADT",
    "facility_name": "ATWOOD-RAWLINS COUNTY CITY-COUNTY",
    "faa_ident": "ADT",
    "type": "AIRPORT",
    "region": "ACE",
    "state_full": "KANSAS",
    "county": "RAWLINS",
    "city": "ATWOOD",
    "ownership": "PU",
    "use": "PU",
    "manager": "LARRY POTTS",
    "manager_phone": "205-373-0446",
    "latitude": "33-06-24.2800N",
    "longitude": "088-11-49.8300W"
}
```

### ✈️ Aviation Service

| Method   | Endpoint | Description                                                                               |
| -------- | -------- | ----------------------------------------------------------------------------------------- |
| **POST** | `/sync`  | Sync incomplete (`PENDING`) airports with AviationAPI. Returns success and failed counts. |

### 🌦️ Weather Service

| Method  | Endpoint                 | Description                                    |
| ------- | ------------------------ | ---------------------------------------------- |
| **GET** | `/weather?city=ASHVILLE` | Get current weather for a city from WeatherAPI |

### 🌍 Airport Weather Service

| Method  | Endpoint                                                          | Description                                                |
| ------- | ----------------------------------------------------------------- | ---------------------------------------------------------- |
| **GET** | `/airport-weather?icao=KADT&facilityName=washington&page=1&pageSize=10` | Get airport data combined with current weather (paginated) |

---

## 🧠 Data Flow Overview

1. Client requests airport data<br>
→ Service first checks Redis cache.<br>
→ If not found, queries PostgreSQL.<br>
→ If still not found, fetches from AviationAPI and stores in both cache + database.

2. Weather requests<br>
→ Directly calls WeatherAPI, cached for short-term reuse.

3. Scheduler<br>
→ Periodically syncs airports with status = `"PENDING"` from the Aviation API.

## 🧪 Testing Tips

1. Run unit tests with coverage:
```bash
go test ./internal/... -cover
```
2. Generate a detailed coverage report:
```bash
go test ./internal/... -coverprofile="coverage.out"
```
3. Visualize test coverage in your browser:
```bash
go tool cover -html="coverage.out"
```

---

## 🧰 Postman Collection
A ready-to-use Postman collection is available in this repository for easier testing and development.

---

## 🧾 License & Credits

* Airport data powered by [AviationAPI](https://www.aviationapi.com/)
* Weather data powered by [WeatherAPI](https://www.weatherapi.com/)