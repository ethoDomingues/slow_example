module github.com/ethodomingues/slow_example/auth

go 1.19

replace github.com/ethodomingues/slow => ../../slow

replace github.com/ethodomingues/slow_example => ../

replace github.com/ethodomingues/slow_example/cdn => ../cdn

replace github.com/ethodomingues/slow_example/api => ../api

require (
	github.com/ethodomingues/slow v0.0.0-20221105170026-0e099cdf7ddf
	github.com/ethodomingues/slow_example v0.0.0-00010101000000-000000000000
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.15 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/exp v0.0.0-20221111204811-129d8d6c17ab // indirect
	gorm.io/driver/sqlite v1.4.3 // indirect
	gorm.io/gorm v1.24.0 // indirect
)
