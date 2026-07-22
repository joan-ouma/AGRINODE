# AGRINODE

AGRINODE is a full-stack, event-driven telemetry platform designed for agricultural IoT. It features a complete closed-loop automation architecture capable of ingesting real-time sensor data, detecting environmental anomalies, and autonomously issuing mitigation commands to edge devices.

##  Architecture

The system is built using a modern microservices architecture, entirely orchestrated via Docker Compose:

###  Data Ingestion & Storage
- **Simulator**: Emulates an ESP8266 microcontroller publishing temperature, humidity, and soil moisture data.
- **MQTT Broker (Mosquitto)**: Handles lightweight M2M communication with edge nodes.
- **Ingestion Service (Go)**: Subscribes to MQTT topics and securely forwards high-volume raw telemetry to Kafka.
- **Kafka / Zookeeper**: The central nervous system, handling high-throughput asynchronous message brokering.
- **Storage Service (Go)**: Consumes the raw data stream from Kafka and persists it to the database.

###  Analytics & Automation
- **PostgreSQL**: The primary operational database. It utilizes SQL triggers and views to compute real-time daily averages and detect critical anomalies (e.g., sudden temperature spikes or soil moisture drops).
- **Event Publisher Service (Go)**: Listens to PostgreSQL database triggers (`anomaly_channel`) and broadcasts detected anomalies back into a Kafka topic.
- **Command Service (Go)**: The autonomous decision engine. It consumes anomaly events from Kafka and fires specific mitigation commands (like `cooling_on` or `pump_on`) back to the specific edge node via MQTT.
- **MongoDB & Analytics Service**: Stores and serves long-term historical data logs for the dashboard.

###  Frontend Dashboard
- **React + Vite**: A modern, glassmorphic UI that provides a mobile-responsive dashboard.
- **Nginx API Gateway**: Routes frontend requests seamlessly to the appropriate backend Go APIs.
- Features real-time charting (Recharts), live device tracking, and historical data logs.

## Getting Started

Ensure you have Docker and Docker Compose installed on your system.

### 1. Build and Run the Stack
```bash
docker compose up -d --build
```
This will spin up all 16 containers including the databases, message brokers, microservices, and the frontend web server.

### 2. Access the Dashboard
Navigate to `http://localhost` in your browser. The dashboard is fully mobile-responsive and will display live incoming telemetry data immediately.

### 3. Monitor the Autonomous Loop
You can watch the system detect anomalies and issue commands by checking the command service logs:
```bash
docker compose logs -f agrinode-command
```

## Tech Stack
- **Backend**: Go (Golang)
- **Frontend**: React.js, Vite, Lucide React, Recharts
- **Databases**: PostgreSQL, MongoDB
- **Brokers**: Apache Kafka, Eclipse Mosquitto (MQTT)
- **Infrastructure**: Docker, Docker Compose, Nginx
