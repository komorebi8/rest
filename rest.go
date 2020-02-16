package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
)

type Rest struct {
	Engine   *gin.Engine
	DB       *gorm.DB
	models   map[string]*Model
	BathPath string
}

func New(e *gin.Engine, db *gorm.DB) *Rest {
	return &Rest{
		Engine:   e,
		DB:       db,
		BathPath: "/api",
		models:   make(map[string]*Model),  // name of model and model
	}
}

func (r *Rest) AddModel(instance interface{}) error {
	t := reflect.TypeOf(instance)
	if t.Kind() == reflect.Struct {
		r.DB.AutoMigrate(instance)
		r.models[strings.ToLower(t.Name())] = NewModel(instance)
	} else {
		return errors.New("model must be a struct")
	}
	return nil
}

func (r *Rest) ForModel(instance interface{}) *Model {
	t := reflect.TypeOf(instance)
	return r.models[strings.ToLower(t.Name())]
}

func (r *Rest) Run(addr ...string) (err error){
	r.Engine.GET("/", func(c *gin.Context) {
		links := gin.H{}
		for name := range r.models {
			links[name] = gin.H{
				"href" : c.Request.Host + r.BathPath + "/" + name,
			}
		}
		c.JSON(200, gin.H{
			"_links" : links,
		})
	})
	r.Engine.GET(r.BathPath + "/:model", func(c *gin.Context) {
		// todo paging and sorting
		name := c.Param("model")
		r.models[name].GetModelFunc(r, c)
	})
	r.Engine.GET(r.BathPath + "/:model/:id", func(c *gin.Context) {
		name := c.Param("model")
		r.models[name].GetModelIDFunc(r, c)
	})
	r.Engine.POST(r.BathPath + "/:model", func(c *gin.Context) {
		name := c.Param("model")
		r.models[name].PostModelFunc(r, c)
	})
	r.Engine.DELETE(r.BathPath + "/:model/:id", func(c *gin.Context) {
		name := c.Param("model")
		r.models[name].DeleteModelIDFunc(r, c)
	})
	r.Engine.PUT(r.BathPath + "/:model/:id", func(c *gin.Context) {
		name := c.Param("model")
		r.models[name].PutModelIDFunc(r, c)
	})
	return r.Engine.Run(addr...)
}
