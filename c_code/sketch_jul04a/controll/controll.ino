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
        Serial.println("Connecting to WiFi...");
        delay(1000);
    }
    Serial.println("Connected to WiFi");

    // Login olup JWT token'ını al
    if (!login("test", "123")) {
        Serial.println("Login failed!");
        // Hata durumunda cihazı resetleyerek veya başka bir kurtarma mekanizması kullanarak yeniden denemek isteyebilirsiniz.
        ESP.reset(); // Cihazı resetleme örneği
    }
}

void loop() {
    checkDoor();
    checkMovement();
}

void checkDoor() {
    int doorState = digitalRead(doorPin);
    if (doorState == HIGH && !doorFlag) {
        Serial.println("Door opened");
        if (sendRequest("open_door", "{\"door\":\"Server Otagyn Gapysy Acyldy\"}", "1")) {
            activateAlarm();
            deactivateAlarm();
            doorFlag = true;
        }
    } else if (doorState == LOW && doorFlag) {
        if (sendRequest("open_door", "{\"door\":\"Server Otagyn Gapysy Yapyldy\"}", "1")) {
            doorFlag = false;
            Serial.println("Door closed");
        }
    }
    delay(100);
}

void checkMovement() {
    int pirPinState = digitalRead(pirPin);
    if (pirPinState == LOW && !pirFlag) {
        Serial.println("Movement detected");
        if (sendRequest("movement_alert", "{\"pir\":\"Hereket Bar\"}", "1")) {
            activateAlarm();
            deactivateAlarm();
            pirFlag = true;
        }
    } else if (pirPinState == HIGH && pirFlag) {
        if (sendRequest("movement_alert", "{\"pir\":\"Hereket yok\"}", "1")) {
            pirFlag = false;
            Serial.println("Movement stopped");
        }
    }
    delay(100);
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
