# GO-REST

This is my first try on Golang.
There is a struct to represent the data and another struct
to store error messages before sending them upstream

Complete CRUD using [gorm.io](htttps://gorm.io) with postgres.

## How to run?
Just use `go run ./main.go`

or build and run the executable

## API endpoints
- /
- /notes `POST` and `GET`
- /notes/{id} `GET`, `PUT`, and `DELETE`