package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kunnpuu/rest"
	"os"
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
	e := gin.Default()
	os.Remove("test.db")
	db, _ := gorm.Open("sqlite3", "test.db")
	defer db.Close()
	r := rest.New(e, db)
	r.AddModel(Product{})
	r.AddModel(Customer{})
	r.ForModel(Customer{}).GetModelFunc = func(r *rest.Rest, c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "customer",
		})
	}
	r.Run()
}
