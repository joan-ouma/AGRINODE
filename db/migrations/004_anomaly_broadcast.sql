-- STEP 1: Create the Broadcast Function
CREATE OR REPLACE FUNCTION notify_anomaly_event()
RETURNS TRIGGER AS $$
DECLARE
    -- We need a variable to hold our data package before we send it
    payload JSON;
BEGIN
    -- Take the brand new anomaly row (NEW) and automatically convert it into a JSON object.
    -- This is incredibly powerful because our Go service and Kafka both natively speak JSON!
    payload = row_to_json(NEW);

    -- Broadcast the JSON payload over a specific channel named 'anomaly_channel'
    -- Think of 'anomaly_channel' as a radio frequency that our Go service will tune into.
    PERFORM pg_notify('anomaly_channel', payload::text);

    -- Let the database finish the insert
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- STEP 2: Arm the Broadcaster
-- Remove the trigger if it already exists to avoid duplicates
DROP TRIGGER IF EXISTS trigger_notify_anomaly ON anomalies;

-- Attach the broadcaster to the anomalies table
-- It will fire AFTER a new warning is successfully saved to the ledger
CREATE TRIGGER trigger_notify_anomaly
AFTER INSERT ON anomalies
FOR EACH ROW
EXECUTE FUNCTION notify_anomaly_event();
