# go-backend

- Framework : Echo
- ORM : GORM
- MongoDB : mongo driver
- GUI based API testing : swagger

### Setup
- Install `air` from [repo](https://github.com/air-verse/air)
- Run `go mod tidy`
- Run `air`

### Architecture
```
.
├── infrastructure/rest/
│   ├── rest
│   └── repository
├── domain
└── application/
    ├── adapter
    └── abi
```
