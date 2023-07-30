const bleno = require('bleno');
const wifiConfig = require('./wifiConfig');
const wifiCheck = require('./wifiStatus');

const myCameraTowerServiceUuid = 'fb0af608-c3ad-41bb-9aba-6d8185f45de7';
const writeCharacteristicUuid = '0cb87266-9c1e-4e8b-a317-b742364e03b4';
const notifyCharacteristicUuid = '53b1b0ed-b315-4665-8b5a-315fc594d84f';

const wifiCredentials = {
  ssid: '',
  psk: '',
  key_mgmt: 'WPA-PSK',
};

// Function to update the Notify characteristic
let updateNotifyCharacteristic;

const writeCharacteristic = new bleno.Characteristic({
  uuid: writeCharacteristicUuid,
  properties: ['write'],
  onWriteRequest: (data, offset, withoutResponse, callback) => {
    // Handle the data received from the web app
    const receivedMessage = data.toString('utf-8').split(':');
    wifiCredentials.ssid = receivedMessage[0];
    wifiCredentials.psk = receivedMessage[1];
    console.log('Received message:', wifiCredentials);

    // Process the receivedMessage (e.g., configure Wi-Fi settings)
    // and restart pi for settings to take hold and connect
    wifiConfig.writeWifiCredentials(wifiCredentials);

    // Send a response back to the web app (optional)
    const response = 'Message received by Raspberry Pi';
    callback(bleno.Characteristic.RESULT_SUCCESS, Buffer.from(response, 'utf-8'));
  },
});

const notifyCharacteristic = new bleno.Characteristic({
  uuid: notifyCharacteristicUuid,
  properties: ['notify'],
  onSubscribe: (maxValueSize, updateValueCallback) => {
    console.log('Central subscribed to notifications');
    updateNotifyCharacteristic = updateValueCallback;
  },
  onUnsubscribe: () => {
    console.log('Central unsubscribed from notifications');
    updateNotifyCharacteristic = null;
  },
});

// Function to notify the central (web app)
const notifyCentral = (isConnected) => {
  if (updateNotifyCharacteristic) {
    const data = Buffer.from([isConnected ? 1 : 0]);
    updateNotifyCharacteristic(data);
  }
};

setInterval( () => {
  // Check wifi status 
  wifiCheck.wifiCheck((result) => {
    if (result) {
      console.log('Device connected to WiFi');
      notifyCentral(true); // Notify the central (web app) that the device is connected
    }
    // console.log('Result:', result);
  });
}, 3000);


const myCameraTowerService = new bleno.PrimaryService({
  uuid: myCameraTowerServiceUuid,
  characteristics: [writeCharacteristic, notifyCharacteristic], //helloCharacteristicUuid
});

bleno.on('stateChange', (state) => {
  if (state === 'poweredOn') {
    bleno.startAdvertising('MyCameraTower', [myCameraTowerServiceUuid]);
  } else {
    bleno.stopAdvertising();
  }
});

bleno.on('advertisingStart', (error) => {
  if (!error) {
    console.log('Advertising as MyCameraTower');
    bleno.setServices([myCameraTowerService]);
  } else {
    console.error('Error starting advertising:', error);
  }
});

bleno.on('accept', (clientAddress) => {
  console.log('Accepted connection from:', clientAddress);
});

bleno.on('disconnect', (clientAddress) => {
  console.log('Disconnected from:', clientAddress);
});
