#include <TinyWireS.h>
#include <EEPROM.h>

#define PIN_SWITCH PB3

#define ADDRESS_I2C_SLAVE 0x04
#define ADDRESS_TAMPER_ADDR 0                       

#define VALUE_TAMPERED 20
#define CMD_RESET_TAMPER_STATE 40
#define CMD_GET_TAMPER_STATE 41

byte cachedTamperState = 0;

void receiveCommand() {                  
  byte command = TinyWireS.receive();
  if (command == CMD_RESET_TAMPER_STATE) {
    resetTamperState();
  } else if (command == CMD_GET_TAMPER_STATE) {
    TinyWireS.send(getTamperState());
  }
}

void setTampered() {
  // Do nothing if already tampered
  if (cachedTamperState == VALUE_TAMPERED) {
    return;
  }

  // First update cache
  cachedTamperState = VALUE_TAMPERED;

  // Then update EEPROM
  byte tamperAddr = EEPROM.read(ADDRESS_TAMPER_ADDR);
  EEPROM.write(tamperAddr, VALUE_TAMPERED);
}

byte getTamperState() {
  // Use cached value to avoid EEPROM read
  return cachedTamperState;
}

void resetTamperState() {  
  // Assign untampered value to the new tamper address
  byte newUntamperedValue = 0;
  do {
    newUntamperedValue = random(255);
  } while (newUntamperedValue == VALUE_TAMPERED);

  cachedTamperState = newUntamperedValue;

  // Assign a new tamper address
  byte newTamperAddr = 0;
  do {
    newTamperAddr = random(EEPROM.length());
  } while (newTamperAddr == ADDRESS_TAMPER_ADDR);

  EEPROM.write(ADDRESS_TAMPER_ADDR, newTamperAddr);    
  EEPROM.write(newTamperAddr, newUntamperedValue);
}

void setup() {
  pinMode(PIN_SWITCH, INPUT);

  initCachedTamperState();

  TinyWireS.begin(ADDRESS_I2C_SLAVE);
  TinyWireS.onReceive(receiveCommand);
}

void initCachedTamperState() {
  byte tamperAddr = EEPROM.read(ADDRESS_TAMPER_ADDR);
  cachedTamperState = EEPROM.read(tamperAddr);
}

void loop() {
  if (digitalRead(PIN_SWITCH) == HIGH) {
    setTampered();
  }
}
