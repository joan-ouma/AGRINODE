-- STEP 1: The Daily Summary (Virtual Table)
-- This creates a view that automatically calculates our daily averages so our Go app doesn't have to do the math.
CREATE OR REPLACE VIEW daily_node_averages AS
SELECT
    n.name AS node_name,
    n.zone,
    DATE(t.created_at) AS reading_date,

    -- Calculate averages and round them so the numbers look clean
    ROUND(AVG(t.temperature)::numeric, 2) AS avg_temp,
    ROUND(AVG(t.humidity)::numeric, 2) AS avg_humidity,
    ROUND(AVG(t.soil_moisture)::numeric, 0) AS avg_moisture,

    -- Keep a running total of how many readings we got today
    COUNT(t.id) AS total_daily_readings
FROM
    nodes n
    JOIN
    telemetry t ON n.id = t.node_id
GROUP BY
    n.name, n.zone, DATE(t.created_at)
ORDER BY
    reading_date DESC;


-- STEP 2: The Security Guard (Anomaly Detection Logic)
-- This function runs the math to see if a brand new reading is dangerously abnormal.
CREATE OR REPLACE FUNCTION check_telemetry_anomaly()
RETURNS TRIGGER AS $$
DECLARE
    -- Create temporary containers to hold our moving averages
avg_temp REAL;
    avg_hum REAL;
    avg_moist REAL;
BEGIN
    -- Look at the last 10 readings for this specific sensor and calculate the average
SELECT AVG(temperature), AVG(humidity), AVG(soil_moisture)
INTO avg_temp, avg_hum, avg_moist
FROM (
         SELECT temperature, humidity, soil_moisture
         FROM telemetry
         WHERE node_id = NEW.node_id
         ORDER BY created_at DESC
             LIMIT 10
     ) AS last_10;

-- Only run the security checks if we actually have enough past data to form an average
IF avg_temp IS NOT NULL THEN

        -- CHECK 1: Is it suddenly way too hot? (5 degrees above normal)
        IF NEW.temperature > (avg_temp + 5.0) THEN
            INSERT INTO anomalies (node_id, anomaly_type, current_value, moving_average)
            VALUES (NEW.node_id, 'High Temperature Spike', NEW.temperature, avg_temp);
END IF;

        -- CHECK 2: Did the air humidity suddenly crash? (10% below normal)
        IF NEW.humidity < (avg_hum - 10.0) THEN
             INSERT INTO anomalies (node_id, anomaly_type, current_value, moving_average)
             VALUES (NEW.node_id, 'Sudden Humidity Drop', NEW.humidity, avg_hum);
END IF;

        -- CHECK 3: Is the soil suddenly drying out? (100 points below normal)
        -- This covers irrigation failures or intense droughts!
        IF NEW.soil_moisture < (avg_moist - 100) THEN
             INSERT INTO anomalies (node_id, anomaly_type, current_value, moving_average)
             VALUES (NEW.node_id, 'Critically Dry Soil', NEW.soil_moisture, avg_moist);
END IF;

END IF;

    -- Let the database finish saving the new reading
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- STEP 3: Turn the Security Guard On
-- First, remove the old trigger if it exists so we don't accidentally make duplicates
DROP TRIGGER IF EXISTS trigger_check_anomaly ON telemetry;

-- Tell the database to automatically run our Security Guard function EVERY time a new row is saved
CREATE TRIGGER trigger_check_anomaly
    AFTER INSERT ON telemetry
    FOR EACH ROW
    EXECUTE FUNCTION check_telemetry_anomaly();