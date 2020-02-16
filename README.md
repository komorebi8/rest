# rest
A tiny Go framework like Spring Data Rest.

It can generate restful api quickly by defining a struct.

For example, writing this code, and you can GET and POST on "localhost:8080/api/product",
also GET, POST, DELETE and PUT on "localhost:8080/api/product/:id".

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

    
    /*
    GET: localhost:8080/
    Response:
    {
        "_links": {
            "customer": {
                "href": "localhost:8080/api/customer"
            },
            "product": {
                "href": "localhost:8080/api/product"
            }
        }
    }
    
    POST: localhost:8080/api/product
    Body:
    {
        "Code": "H234567",
        "Price": 23467
    }
    Response:
    {
        "ID": 1,
        "CreatedAt": "2020-02-06T12:08:24.7030646+08:00",
        "UpdatedAt": "2020-02-06T12:08:24.7030646+08:00",
        "DeletedAt": null,
        "Code": "H234567",
        "Price": 23467
    }
    
    GET: localhost:8080/api/product
    Response:
    {
        "_embedded": {
            "product": [
                {
                    "ID": 1,
                    "CreatedAt": "2020-02-06T12:08:24.7030646+08:00",
                    "UpdatedAt": "2020-02-06T12:08:24.7030646+08:00",
                    "DeletedAt": null,
                    "Code": "H234567",
                    "Price": 23467
                }
            ]
        },
        "_links": {
            "self": {
                "href": "localhost:8080/api/product"
            }
        }
    }
    */

