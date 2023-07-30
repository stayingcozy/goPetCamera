const fs = require('fs');
const { exec } = require('child_process');

function writeWifiCredentials(wifiNetwork) {

// Create the content for the wpa_supplicant.conf file
const wpaConfigContent = `ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1
country=US

network={
    ssid="${wifiNetwork.ssid}"
    psk="${wifiNetwork.psk}"
    key_mgmt=${wifiNetwork.key_mgmt}
}
`; 

// Define the file path for wpa_supplicant.conf
const wpaConfigFilePath = '/etc/wpa_supplicant/wpa_supplicant.conf';

// Write the content to the file
fs.writeFile(wpaConfigFilePath, wpaConfigContent, (err) => {
    if (err) {
    console.error('Error writing wpa_supplicant.conf:', err);
    } else {
    console.log('wpa_supplicant.conf file created/updated successfully!');
    restartWifiInterface()
    }
});
}

function restartWifiInterface() {
  // Reconfigure wifi
    exec('wpa_cli -i wlan0 reconfigure', (error, stdout, stderr) => {
    if (error) {
        console.error(`Error while reconfigure: ${error.message}`);
        return;
    }
    })
};

module.exports = { writeWifiCredentials };
