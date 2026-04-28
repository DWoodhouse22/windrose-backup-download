# Windrose World Save Downloader
Simple go code which downloads the current live world save files for semi-regular backups

## Usage
Download the executable for your platform from the latest release

Create a config file using the example provided (`config.example.json`) in the same directory you downloaded the executable. Replace the values with your FTP credentials and Windrose World ID

Run the executable, by default it will load the file `config.json`. Specify a custom file by providing the file path `./windrose-backup-windows-amd64.exe path/to/my/custom-config.json`

Resulting download will be placed in a zip file with the name of your world ID and the current timestamp