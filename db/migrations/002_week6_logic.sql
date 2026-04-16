-- 1. Create the Daily Averages View
CREATE OR REPLACE VIEW daily_node_averages AS
SELECT 
    n.name AS node_name,
    n.zone,
    DATE(t.created_at) AS reading_date,
    ROUND(AVG(t.temperature)::numeric, 2) AS avg_temp,
    ROUND(AVG(t.humidity)::numeric, 2) AS avg_humidity,
    ROUND(AVG(t.soil_moisture)::numeric, 0) AS avg_moisture,
    COUNT(t.id) AS total_daily_readings
FROM 
    nodes n
JOIN 
    telemetry t ON n.id = t.node_id
GROUP BY 
    n.name, n.zone, DATE(t.created_at)
ORDER BY 
    reading_date DESC;

-- 2. Create the Anomaly Detection Trigger Function
CREATE OR REPLACE FUNCTION check_telemetry_anomaly()
RETURNS TRIGGER AS $$
DECLARE
    avg_temp REAL;
    avg_hum REAL;
BEGIN
    -- Calculate the moving average of the last 10 readings for THIS specific node
    SELECT AVG(temperature), AVG(humidity)
    INTO avg_temp, avg_hum
    FROM (
        SELECT temperature, humidity
        FROM telemetry
        WHERE node_id = NEW.node_id
        ORDER BY created_at DESC
        LIMIT 10
    ) AS last_10;

    -- Only check if we have enough data to form an average
    IF avg_temp IS NOT NULL THEN
        -- Check for a Temperature Spike
        IF NEW.temperature > (avg_temp + 5.0) THEN
            INSERT INTO anomalies (node_id, anomaly_type, current_value, moving_average)
            VALUES (NEW.node_id, 'High Temperature Spike', NEW.temperature, avg_temp);
        END IF;

        -- Check for Sudden Humidity Drop
        IF NEW.humidity < (avg_hum - 10.0) THEN
             INSERT INTO anomalies (node_id, anomaly_type, current_value, moving_average)
             VALUES (NEW.node_id, 'Sudden Humidity Drop', NEW.humidity, avg_hum);
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 3. Arm the Trigger
DROP TRIGGER IF EXISTS trigger_check_anomaly ON telemetry;

CREATE TRIGGER trigger_check_anomaly
AFTER INSERT ON telemetry
FOR EACH ROW
EXECUTE FUNCTION check_telemetry_anomaly();
