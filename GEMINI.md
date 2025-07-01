# Quick linker app
Is an application that shortens urls and create QR codes for them.

## Instructions
- Use golang to build the app
- Use gin framework to build the api
- Use redis to store the urls
- Use qrcode to generate the qr codes
- Use docker to containerize the app
- Use docker compose to run the app


## Arquitecture
Has a rest api to shorten urls and a rest api to generate qr codes.
-Golang pattern organization


## Rules of buisness
1. The user can short and url if is valid
2. The user can generate a qr code for a url
3. The user can get the stats of a url
4. The user can delete a url
5. The user can determine the expiration date or time of a url
6. The app must have ad's

## Structure of the project
quick_linker/
├── front-end/
├── back-end/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── handler/
│   │   │   ├── handler.go
│   │   │   └── handler_test.go
│   │   ├── repository/
│   │   │   ├── repository.go
│   │   │   └── repository_test.go
│   │   ├── service/
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   │   └── utils/
│   │       ├── utils.go
│   │       └── utils_test.go
│   └── go.mod
├── docker-compose.yaml
├── go.mod
├── go.sum
└── README.md
