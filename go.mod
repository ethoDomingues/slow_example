module github.com/ethodomingues/slow_example

go 1.19

replace github.com/ethodomingues/slow => ../slow

require (
	github.com/ethodomingues/slow v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.1.0
	gorm.io/driver/sqlite v1.4.3
	gorm.io/gorm v1.24.0
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.15 // indirect
)
