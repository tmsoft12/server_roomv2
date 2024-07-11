#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266HTTPClient.h>
#include <DHT.h>

const char* ssid = "Osman";
const char* password = "21032020";
const char* serverUrl = "http://192.168.100.191:3000/";
const char* loginUrl = "http://192.168.100.191:3000/login";
const int firePin = 14;
const int dhtPin = 12;

bool fireFlag = false;
String jwtToken; 
DHT dht(dhtPin, DHT11);
float lastTemp = 0.0;
float lastHum = 0.0;

void setup() {
    Serial.begin(115200);
    pinMode(firePin, INPUT);
    dht.begin();
    
    connectWiFi();

    
    if (!login("test", "123")) {
        Serial.println("Login failed!");
        ESP.reset(); 
    }
}

void loop() {
    checkFire();
    checkTempHum();
    delay(100);
}

void connectWiFi() {
    WiFi.begin(ssid, password);
    while (WiFi.status() != WL_CONNECTED) {
        Serial.println("Connecting to WiFi...");
        delay(1000);
    }
    Serial.println("Connected to WiFi");
}

void checkFire() {
    int fireState = digitalRead(firePin);
    if (fireState == HIGH && !fireFlag) {
        Serial.println("Fire detected");
        if (sendRequest("fire_alert", "{\"fire\":\"Fire detected\"}", "1")) {
            fireFlag = true;
        }
    } else if (fireState == LOW && fireFlag) {
        Serial.println("Fire extinguished");
        if (sendRequest("fire_alert", "{\"fire\":\"Fire extinguished\"}", "1")) {
            fireFlag = false;
        }
    }
}

void checkTempHum() {
    float newTemp = dht.readTemperature();
    float newHum = dht.readHumidity();

    if (isnan(newTemp) || isnan(newHum)) {
        Serial.println("Failed to read from DHT sensor!");
        return;
    }

    if (newTemp != lastTemp || newHum != lastHum) {
        Serial.print("Temperature: ");
        Serial.print(newTemp);
        Serial.print(" Humidity: ");
        Serial.println(newHum);

        String jsonPayload = "{\"temp\":\"" + String(newTemp) + "\",\"hum\":\"" + String(newHum) + "\"}";
        if (sendRequest("temp_hum", jsonPayload.c_str(), "1")) {
            lastTemp = newTemp;
            lastHum = newHum;
        }
    }
}

bool sendRequest(const char* endpoint, const char* jsonPayload, const char* id) {
    String url = serverUrl;
    url += endpoint;
    url += "/";
    url += id;

    WiFiClient client;
    HTTPClient http;

    http.begin(client, url);
    http.addHeader("Authorization", "Bearer " + jwtToken);
    http.addHeader("Content-Type", "application/json");

    Serial.print("Sending PUT request to: ");
    Serial.println(url);
    Serial.print("Payload: ");
    Serial.println(jsonPayload);

    int httpResponseCode = http.PUT(jsonPayload);

    if (httpResponseCode > 0) {
        Serial.print("HTTP response code: ");
        Serial.println(httpResponseCode);
        String response = http.getString();
        Serial.println("Response: " + response);
        http.end();
        return true;
    } else {
        Serial.print("Error on sending PUT request: ");
        Serial.println(httpResponseCode);
        String response = http.getString();
        Serial.println("Response: " + response);
        http.end();
        return false;
    }
}

bool login(const char* username, const char* password) {
    WiFiClient client;
    HTTPClient http;

    http.begin(client, loginUrl);
    http.addHeader("Content-Type", "application/json");

    String postData = "{\"username\":\"";
    postData += username;
    postData += "\",\"password\":\"";
    postData += password;
    postData += "\"}";

    int httpResponseCode = http.POST(postData);

    if (httpResponseCode > 0) {
        String response = http.getString();
        Serial.print("HTTP response code: ");
        Serial.println(httpResponseCode);
        Serial.println("Response: " + response);

        // Extract JWT token from response
        int tokenIndex = response.indexOf("\"token\":\"") + 9;
        int tokenEndIndex = response.indexOf("\"", tokenIndex);
        jwtToken = response.substring(tokenIndex, tokenEndIndex);

        Serial.println("JWT Token: " + jwtToken);
        http.end();
        return true;
    } else {
        Serial.print("Error on sending POST request: ");
        Serial.println(httpResponseCode);
        http.end();
        return false;
    }
}
