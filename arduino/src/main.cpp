#include <Adafruit_Sensor.h>
#include <Wire.h>

#define FAN1_PWM_PIN 12     // D6
#define FAN1_SIGNAL_PIN 13  // D7
#define FAN2_PWM_PIN 0      // D3
#define FAN2_SIGNAL_PIN 2   // D4
#define FAN3_PWM_PIN 5      // D1
#define FAN3_SIGNAL_PIN 4   // D2
#define FAN4_PWM_PIN 15     // D8
#define FAN4_SIGNAL_PIN 14  // D5

#define EMERGENCY_SPEED 100

unsigned long lastWatchdogReceived = 0;
const unsigned long WATCHDOG_PERIOD = 30000;

int getFanSpeedRpm(int SIGNAL_PIN) {
  int highTime = pulseIn(SIGNAL_PIN, HIGH);
  int lowTime = pulseIn(SIGNAL_PIN, LOW);
  int period = highTime + lowTime;
  if (period == 0) {
    return 0;
  }
  float freq = 1000000.0 / (float)period;
  return (freq * 60.0) / 2.0;
}

void setFanSpeedPercent(int PWM_PIN, int p) {
  int value = (p / 100.0) * 255;
  analogWriteFreq(25000);
  analogWrite(PWM_PIN, value);
}

void emergencyMode() {
  Serial.println("EMERGENCY MODE!");

  setFanSpeedPercent(FAN1_PWM_PIN, EMERGENCY_SPEED);
  setFanSpeedPercent(FAN2_PWM_PIN, EMERGENCY_SPEED);
  setFanSpeedPercent(FAN3_PWM_PIN, EMERGENCY_SPEED);
  setFanSpeedPercent(FAN4_PWM_PIN, EMERGENCY_SPEED);
}

void resetWatchdog() { lastWatchdogReceived = millis(); }

void checkWatchdog() {
  unsigned long currentTime = millis();
  if (currentTime - lastWatchdogReceived >= WATCHDOG_PERIOD) {
    emergencyMode();
    lastWatchdogReceived = millis();
  }
}

void processCommand(String cmd) {
  String command = cmd.substring(0, cmd.indexOf(' '));

  if (command == "WATCHDOG") {
    resetWatchdog();
    return;
  }

  if (command == "SET") {
    int fan = cmd.substring(cmd.indexOf(' ') + 1, cmd.lastIndexOf(' ')).toInt();
    int pwm = cmd.substring(cmd.lastIndexOf(' ') + 1).toInt();
    int pin;
    switch (fan) {
      case 1:
        pin = FAN1_PWM_PIN;
        break;
      case 2:
        pin = FAN2_PWM_PIN;
        break;
      case 3:
        pin = FAN3_PWM_PIN;
        break;
      case 4:
        pin = FAN4_PWM_PIN;
        break;
      default:
        Serial.println("UNKNOWN FAN!");
        return;
    }

    setFanSpeedPercent(pin, pwm);
    Serial.printf("SET FAN %d TO %d%%\n", fan, pwm);
    return;
  }

  if (command == "GET") {
    int fan = cmd.substring(cmd.lastIndexOf(' ') + 1).toInt();
    int pin;
    switch (fan) {
      case 1:
        pin = FAN1_SIGNAL_PIN;
        break;
      case 2:
        pin = FAN2_SIGNAL_PIN;
        break;
      case 3:
        pin = FAN3_SIGNAL_PIN;
        break;
      case 4:
        pin = FAN4_SIGNAL_PIN;
        break;
      default:
        Serial.println("UNKNOWN FAN!");
        return;
    }

    int rpm = getFanSpeedRpm(pin);
    Serial.println(rpm);
    return;
  }

  Serial.println("UNKNOWN COMMAND!");
}

void setup() {
  Serial.begin(115200);
  Serial.setTimeout(500);
  Serial.print("STARTING...\n");
  pinMode(FAN1_PWM_PIN, OUTPUT);
  pinMode(FAN1_SIGNAL_PIN, INPUT_PULLUP);
  pinMode(FAN2_PWM_PIN, OUTPUT);
  pinMode(FAN2_SIGNAL_PIN, INPUT_PULLUP);
  pinMode(FAN3_PWM_PIN, OUTPUT);
  pinMode(FAN3_SIGNAL_PIN, INPUT_PULLUP);
  pinMode(FAN4_PWM_PIN, OUTPUT);
  pinMode(FAN4_SIGNAL_PIN, INPUT_PULLUP);

  emergencyMode();
}

void loop() {
  checkWatchdog();

  if (Serial.available()) {
    String cmd = Serial.readStringUntil('\n');
    cmd.trim();
    processCommand(cmd);
  }
}