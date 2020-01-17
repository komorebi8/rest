# rest
A tiny Go framework like Spring Data Rest.

It can generate restful api quickly by defining a struct.

For example, writing this code, and you can GET and POST on "localhost:8080/api/product",
also POST, DELETE and PUT on "localhost:8080/api/product/:id".

    package main
    import (
	    "github.com/gin-gonic/gin"
	    "github.com/jinzhu/gorm"
	    _ "github.com/jinzhu/gorm/dialects/sqlite"
	    "github.com/kunnpuu/rest"
    )
    type Product struct {
	    gorm.Model
	    Code string
	    Price uint
    }
    type Customer struct {
	    gorm.Model
	    Name string
	    Age uint
    }
    func main() {
	    r := gin.Default()
	    db, _ := gorm.Open("sqlite3", "test.db")
	    defer db.Close()
	    e := rest.NewEngine(r, db)
	    e.AddModel(Product{})
	    e.AddModel(Customer{})
	    e.Run()
    }

