# yitu

## Installation

### bimg

https://github.com/h2non/bimg#installation

```bash
curl -s https://raw.githubusercontent.com/h2non/bimg/master/preinstall.sh | sudo bash -
```

## Authorization

https://docs.microsoft.com/en-us/onedrive/developer/rest-api/getting-started/graph-oauth?view=odsp-graph-online#code-flow

You need to get access token and refresh token first. You can read the official document from Microsoft above, or follow the steps below.

### Get authorization code with your browser

https://login.microsoftonline.com/common/oauth2/v2.0/authorize?access_type=offline&client_id=4caae01e-515a-490f-bde7-92cff3b895ac&redirect_uri=http%3A%2F%2Flocalhost%3A23333%2F&response_type=code&scope=Files.Read+Files.ReadWrite+Files.Read.All+Files.ReadWrite.All+offline_access

Open this url with your browser. After a few clicks, you will be redirect to `http://localhost:23333/`, the authorization code is in the url params.

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