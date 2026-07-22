#include <WiFi.h>
#include <PubSubClient.h>
#include <DHT.h>
#include <ArduinoJson.h>

// ---------------- CONFIGURATION ----------------
const char* ssid = "YOUR_WIFI_SSID";
const char* password = "YOUR_WIFI_PASSWORD";

// This must be the IP address of your Ubuntu laptop running Docker
// Do NOT use "localhost" or "127.0.0.1" here.
const char* mqtt_server = "192.168.X.X"; 
const int mqtt_port = 1883;
const char* mqtt_user = "agrinode_device";
const char* mqtt_pass = "farm_secret";

const char* pub_topic = "agrinode/telemetry";
const char* sub_topic = "agrinode/commands/1"; // Matches Node 1 in command-svc

// ---------------- PIN DEFINITIONS ----------------
#define DHTPIN 4
#define DHTTYPE DHT11
#define MOISTURE_PIN 34 // ESP32 Analog Pin

// Mitigation Relay Pins
#define PUMP_PIN 18
#define COOLER_PIN 19
#define MISTER_PIN 21

// ---------------- GLOBALS ----------------
WiFiClient espClient;
PubSubClient client(espClient);
DHT dht(DHTPIN, DHTTYPE);

unsigned long lastMsg = 0;

void setup() {
  Serial.begin(115200);
  dht.begin();

  pinMode(PUMP_PIN, OUTPUT);
  pinMode(COOLER_PIN, OUTPUT);
  pinMode(MISTER_PIN, OUTPUT);
  
  // Ensure relays are off by default
  digitalWrite(PUMP_PIN, LOW);
  digitalWrite(COOLER_PIN, LOW);
  digitalWrite(MISTER_PIN, LOW);

  setup_wifi();
  client.setServer(mqtt_server, mqtt_port);
  client.setCallback(mqttCallback);
}

void setup_wifi() {
  delay(10);
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("\nWiFi connected");
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());
}

// ---------------- COMMAND CONSUMER ----------------
void mqttCallback(char* topic, byte* payload, unsigned int length) {
  Serial.print("Command arrived on topic: ");
  Serial.println(topic);

  // Parse the incoming JSON command from your Go backend
  JsonDocument doc; 
  DeserializationError error = deserializeJson(doc, payload, length);

  if (error) {
    Serial.print("deserializeJson() failed: ");
    Serial.println(error.c_str());
    return;
  }

  const char* action = doc["action"];
  int duration_sec = doc["duration_sec"];

  Serial.print("Action: ");
  Serial.println(action);

  // Execute the mitigation command
  if (strcmp(action, "pump_on") == 0) {
    digitalWrite(PUMP_PIN, HIGH);
    delay(duration_sec * 1000); // Simple blocking delay for mitigation
    digitalWrite(PUMP_PIN, LOW);
  } 
  else if (strcmp(action, "cooling_on") == 0) {
    digitalWrite(COOLER_PIN, HIGH);
    delay(duration_sec * 1000);
    digitalWrite(COOLER_PIN, LOW);
  }
  else if (strcmp(action, "misters_on") == 0) {
    digitalWrite(MISTER_PIN, HIGH);
    delay(duration_sec * 1000);
    digitalWrite(MISTER_PIN, LOW);
  }
}

void reconnect() {
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    // Connect as Node 1
    if (client.connect("ESP32_Node_1", mqtt_user, mqtt_pass)) {
      Serial.println("connected");
      // Subscribe to the command channel immediately
      client.subscribe(sub_topic);
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      delay(5000);
    }
  }
}

// ---------------- TELEMETRY PUBLISHER ----------------
void loop() {
  if (!client.connected()) {
    reconnect();
  }
  client.loop();

  unsigned long now = millis();
  // Publish telemetry every 5 seconds
  if (now - lastMsg > 5000) {
    lastMsg = now;

    float h = dht.readHumidity();
    float t = dht.readTemperature();

    // The ESP32 ADC reads 0-4095. 
    // Dry soil = ~4095, Wet soil = ~1000. We map this to a 0-100% scale.
    int raw_moisture = analogRead(MOISTURE_PIN);
    int moisture_percent = map(raw_moisture, 4095, 1000, 0, 100);
    
    // Constrain to 0-100 just in case it drifts
    if (moisture_percent < 0) moisture_percent = 0;
    if (moisture_percent > 100) moisture_percent = 100;

    // Check if any DHT reads failed
    if (isnan(h) || isnan(t)) {
      Serial.println("Failed to read from DHT sensor!");
      return;
    }

    // Format the payload exactly as your Go Ingestion Service expects
    JsonDocument doc;
    doc["temperature"] = t;
    doc["humidity"] = h;
    doc["soilMoisture"] = moisture_percent;

    char jsonBuffer[512];
    serializeJson(doc, jsonBuffer);

    Serial.print("Publishing telemetry: ");
    Serial.println(jsonBuffer);
    
    client.publish(pub_topic, jsonBuffer);
  }
}
