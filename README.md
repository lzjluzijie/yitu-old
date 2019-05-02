# yitu

## Installation

### NodeJS

```bash
curl -sL https://rpm.nodesource.com/setup_10.x | sudo bash -
curl --silent --location https://dl.yarnpkg.com/rpm/yarn.repo | sudo tee /etc/yum.repos.d/yarn.repo
sudo yum install -y nodejs yarn
```

### bimg

https://github.com/h2non/bimg#installation

```bash
curl -s https://raw.githubusercontent.com/h2non/bimg/master/preinstall.sh | sudo bash -
```

## Authorization

https://docs.microsoft.com/en-us/onedrive/developer/rest-api/getting-started/graph-oauth?view=odsp-graph-online#code-flow

You need to get access token and refresh token first. 

### Get authorization code with your browser

https://login.microsoftonline.com/common/oauth2/v2.0/authorize?access_type=offline&client_id=4caae01e-515a-490f-bde7-92cff3b895ac&redirect_uri=http%3A%2F%2Flocalhost%3A23333%2F&response_type=code&scope=Files.Read+Files.ReadWrite+Files.Read.All+Files.ReadWrite.All+offline_access

Open this url with your browser. After a few clicks, you will be redirect to `http://localhost:23333/`, the authorization code is in the url params.


### Redeem the code for tokens

```http request
POST https://login.microsoftonline.com/common/oauth2/v2.0/token
Content-Type: application/x-www-form-urlencoded

client_id=4caae01e-515a-490f-bde7-92cff3b895ac&redirect_uri=http://localhost:23333/&client_secret=qohmO45%%-jtxUVCAGP372{
&code={code}&grant_type=authorization_code
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