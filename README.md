# yitu

## Installation

### bimg

https://github.com/h2non/bimg#installation

```bash
curl -s https://raw.githubusercontent.com/h2non/bimg/master/preinstall.sh | sudo bash -
```

## Config

yitu.json

```json
{
  "Database": {
    "Driver": "sqlite3",
    "Source": "./yitu.db?parseTime=true&loc=Local"
  },
  "OneDrive": {
    "ClientID": "4caae01e-515a-490f-bde7-92cff3b895ac",
    "ClientSecret": "qohmO45%%-jtxUVCAGP372{",
    "AccessToken": "AccessToken",
    "RefreshToken": "RefreshToken"
  }
}
```

## Backup

### sqlite
```
0 * * * * cd ~/yitu/scripts/backup && bash backup.sh
```