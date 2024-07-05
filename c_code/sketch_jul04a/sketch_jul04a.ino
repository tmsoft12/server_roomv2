#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>

const char* ssid = "Osman";
const char* password = "21032020";
const char* serverName = "http://192.168.100.186:3000/login"; 

WiFiClient client;
void setup() {
  Serial.begin(115200);
  
  // Connect to Wi-Fi
  WiFi.begin(ssid, password);
  
  while (WiFi.status() != WL_CONNECTED) {
    delay(2000);
    Serial.println("Connecting to WiFi...");
  }
  
  Serial.println("Connected to WiFi");

  // Send POST request to Go server
  if(WiFi.status() == WL_CONNECTED) {
    HTTPClient http;
    
    // Use client object when calling begin
    http.begin(client, serverName);
    http.addHeader("Content-Type", "application/json");

    String postData = "{\"username\":\"test\",\"password\":\"123\"}";

    int httpResponseCode = http.POST(postData);

    if (httpResponseCode > 0) {
      String response = http.getString();
      Serial.println(httpResponseCode);
      Serial.println(response);
      
      // Extract JWT token from response
      int tokenIndex = response.indexOf("\"token\":\"") + 9;
      int tokenEndIndex = response.indexOf("\"", tokenIndex);
      String jwtToken = response.substring(tokenIndex, tokenEndIndex);
      
      Serial.println("JWT Token: " + jwtToken);
    } else {
      Serial.print("Error on sending POST: ");
      Serial.println(httpResponseCode);
    }
    
    http.end();
  }
}

void loop() {
  // Perform any other tasks in loop if needed
}
