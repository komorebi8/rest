package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"testing"
)

type Product struct {
	gorm.Model
	Code string
	Price uint
}

func TestNewEngine(t *testing.T) {



	r := gin.Default()

	db, _ := gorm.Open("sqlite3", "test.db")
	defer db.Close()

	p1 := &Product{
		Model: gorm.Model{},
		Code:  "H123",
		Price: 123,
	}
	p2 := &Product{
		Model: gorm.Model{},
		Code:  "H134",
		Price: 234,
	}
	db.Create(p1).Create(p2)


	e := NewEngine(r, db)
	e.AddModel(Product{})
	e.Run()

}
