package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"testing"
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

func TestNewEngine(t *testing.T) {
	e := gin.Default()
	os.Remove("test.db")
	db, _ := gorm.Open("sqlite3", "test.db")
	defer db.Close()
	r := New(e, db)
	r.AddModel(Product{})
	r.AddModel(Customer{})
	p1 := &Product{
		Model: gorm.Model{},
		Code:  "H1234",
		Price: 123,
	}
	p2 := &Product{
		Model: gorm.Model{},
		Code:  "H2345",
		Price: 234,
	}
	p3 := &Product{
		Model: gorm.Model{},
		Code:  "H3456",
		Price: 345,
	}
	db.Create(p1).Create(p2).Create(p3)
	r.Run()
}
