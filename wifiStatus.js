const { exec } = require('child_process');

function wifiCheck(callback) {
  // Check wifi status
  const command = 'iwconfig wlan0 | grep "ESSID"';

  exec(command, (error, stdout, stderr) => {
    if (error) {
      console.error(`Error while executing command: ${error.message}`);
      callback(0); // Pass the result to the callback
      return;
    }

    const output = stdout.toString();
    const isConnected = !output.includes('ESSID:off/any');

    if (isConnected) {
    //   console.log('WiFi is connected.');
      callback(1); // WiFi is connected, pass 1 to the callback
    } else {
    //   console.log('WiFi is not connected.');
      callback(0); // WiFi is not connected, pass 0 to the callback
    }
  });
}

module.exports = { wifiCheck };