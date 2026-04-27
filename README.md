# Windrose World Save Downloader
Simple go code which downloads the current live world save files for semi-regular backups

## Requirements
Go runtime for your system:  
[Download](https://go.dev/dl/)

## Usage
- Create a `.env` file at the root of the project and add the following variables.

(Assuming Windrose is hosted on Nitrado)
Retrieve the FTP credentials from your Nitrado server WebInterface

FTP_USER="USERNAME"  
FTP_PASS="PASSWORD"  
FTP_HOST="HOSTNAME:PORT"  
WORLD_ID="WORLD ID"  
WORLD_SAVE_LOCATION="windrose/R5/Saved/SaveProfiles/Default/RocksDB/0.10.0/Worlds"  

With that set, open a terminal in the root and call `go run main.go`