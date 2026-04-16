-- 1. Create the hardware nodes table
CREATE TABLE IF NOT EXISTS nodes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    zone VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active'
);

-- 2. Add the foreign key to the existing telemetry table
ALTER TABLE telemetry ADD COLUMN IF NOT EXISTS node_id INTEGER REFERENCES nodes(id);

-- 3. Create the anomalies ledger
CREATE TABLE IF NOT EXISTS anomalies (
    id SERIAL PRIMARY KEY,
    node_id INTEGER REFERENCES nodes(id),
    anomaly_type VARCHAR(50) NOT NULL,
    current_value REAL NOT NULL,
    moving_average REAL NOT NULL,
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Insert a dummy node so we have data to test with
INSERT INTO nodes (name, zone) VALUES ('ESP8266-Alpha', 'Zone 1') ON CONFLICT DO NOTHING;
