# CamStream
Go webcam streaming server

## How to use CamStream

1) Pull this repo and run `go build` to compile
2) Visit a web browser and enter the local IP of your host device and the port you have chosen

EG: `http://192.168.1.10:3000/login`

3) Log in using the password you set
4) You will be redirected to the stream. The authentication token rotates every time the server is started

## Config

The config is set in a file called `config.json`

```js
{
  "port": 3000,
   "password": "S3curePassw0rd"
}
```
