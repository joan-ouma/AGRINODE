-- 1. Create the cold-storage archive table
CREATE TABLE IF NOT EXISTS telemetry_archive (
    id INTEGER PRIMARY KEY,
    node_id INTEGER REFERENCES nodes(id),
    temperature REAL,
    humidity REAL,
    soil_moisture INTEGER,
    created_at TIMESTAMP,
    archived_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create the Archiving Procedure
CREATE OR REPLACE PROCEDURE archive_old_telemetry()
LANGUAGE plpgsql
AS $$
BEGIN
    -- Copy data older than 30 days into the archive table
    INSERT INTO telemetry_archive (id, node_id, temperature, humidity, soil_moisture, created_at)
    SELECT id, node_id, temperature, humidity, soil_moisture, created_at
    FROM telemetry
    WHERE created_at < NOW() - INTERVAL '30 days';

    -- Delete the copied data from the live table
    DELETE FROM telemetry
    WHERE created_at < NOW() - INTERVAL '30 days';

    -- Emit a notice to the server logs
    RAISE NOTICE 'Old telemetry data successfully archived.';
END;
$$;
