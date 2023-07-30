const sudo = require('sudo-prompt');
const path = require('path');
const os = require('os');
const { wifiCheck } = require('./wifiStatus');

// Main script
wifiCheck((result) => {

    goPath = path.join(getUserHome(), 'goPetCamera');

    if (result === 1) {
    // WiFi is connected, run the Go program (main)
    const { spawn } = require('child_process');
    const mainPath = path.join(goPath, 'main');

    // Change the current working directory to where the 'main' executable is located
    process.chdir(path.dirname(goPath));

    const mainProcess = spawn(mainPath, [], {
        stdio: 'inherit',
    });

    mainProcess.on('close', (code) => {
        console.log(`Main process exited with code ${code}`);
    });
    } else {
    // WiFi is not connected, run the Node.js BLE script with sudo
    const options = {
        name: 'MyApp', // Provide your app name or identifier here
    };
    blepi_path = path.join(__dirname,"bleno_connect.js");
    sudo.exec('node '+blepi_path, options, (error, stdout, stderr) => {
        if (error) {
        console.error(`Error while executing BLE script with sudo: ${error.message}`);
        return;
        }
        console.log(`BLE script output: ${stdout}`);
        // When the BLE script exits, start the main Go program
        const { spawn } = require('child_process');
        const mainPath = path.join(goPath, 'main');

        // Change the current working directory to where the 'main' executable is located
        process.chdir(path.dirname(goPath));

        const mainProcess = spawn(mainPath, [], {
        stdio: 'inherit',
        });

        mainProcess.on('close', (code) => {
        console.log(`Main process exited with code ${code}`);
        });
    });
    }

});

function getUserHome() {
    return os.homedir();
}
