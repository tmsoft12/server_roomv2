#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266HTTPClient.h>

const char* ssid = "Osman";
const char* password = "21032020";
const char* serverUrl = "http://192.168.100.191:3000/"; // Sunucunuzun adresi
const char* loginUrl = "http://192.168.100.191:3000/login"; // Login endpoint

const int doorPin = 14;
const int pirPin = 12;
const int buzzerPin = 2;

bool doorFlag = false;
bool pirFlag = false;

String jwtToken; // JWT token'ı saklamak için

void setup() {
    Serial.begin(115200);

    pinMode(doorPin, INPUT);
    pinMode(buzzerPin, OUTPUT);
    pinMode(pirPin, INPUT);

    WiFi.begin(ssid, password);
    while (WiFi.status() != WL_CONNECTED) {
        delay(1000);
        Serial.println("Connecting to WiFi...");
    }
    Serial.println("Connected to WiFi");

    // Login olup JWT token'ını al
    if (!login("test", "123")) {
        Serial.println("Login failed!");
        while (true); // Login başarısızsa, işlemi durdur
    }
}

void loop() {
    checkDoor();
    checkMovement();
}

void checkDoor() {
    int doorState = digitalRead(doorPin);
    if (doorState == LOW && !doorFlag) {
        Serial.println("Door opened");
        sendRequest("open_door", "{\"door\":\"Server Otagyn Gapysy Acyldy\"}", "1"); // id parametresi ile çağrı
        activateAlarm();
        deactivateAlarm();
        doorFlag = true;
    } else if (doorState == HIGH && doorFlag) {
        sendRequest("open_door", "{\"door\":\"Server Otagyn Gapysy Yapyldy\"}", "1"); // id parametresi ile çağrı
        doorFlag = false;
        Serial.println("Door closed");
    }
    delay(100);
}

void checkMovement() {
    int pirPinState = digitalRead(pirPin);
    if (pirPinState == LOW && !pirFlag) {
        Serial.println("Movement detected");
        sendRequest("movement_alert", "{\"movement\":\"detected\"}", "1"); // id parametresi ile çağrı
        activateAlarm();
        deactivateAlarm();
        pirFlag = true;
    } else if (pirPinState == HIGH && pirFlag) {
        sendRequest("movement_alert", "{\"movement\":\"stopped\"}", "1"); // id parametresi ile çağrı
        pirFlag = false;
        Serial.println("Movement stopped");
    }
    delay(100);
}

void sendRequest(const char* endpoint, const char* jsonPayload, const char* id) {
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
    } else {
        Serial.print("Error on sending PUT request: ");
        Serial.println(httpResponseCode);
        String response = http.getString();
        Serial.println("Response: " + response);
    }

    http.end();
}

void activateAlarm() {
    tone(buzzerPin, 1000);
}

void deactivateAlarm() {
    noTone(buzzerPin);
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
        Serial.println(httpResponseCode);
        Serial.println(response);
        
        // Yanıttan JWT token'ı çıkarma
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
