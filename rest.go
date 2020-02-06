package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Rest struct {
	Engine   *gin.Engine
	DB       *gorm.DB
	models   map[string]*Model
	BathPath string
}

type TweakFunc func(r *Rest, c *gin.Context)

func New(e *gin.Engine, db *gorm.DB) *Rest {
	return &Rest{
		Engine:   e,
		DB:       db,
		BathPath: "/api",
		models:   make(map[string]*Model),  // name of model and model
	}
}

func (r *Rest) ForModel(model interface{}) *Model {
	t := reflect.TypeOf(model)
	return r.models[strings.ToLower(t.Name())]
}

func (r *Rest) AddModel(model interface{})  {
	r.DB.AutoMigrate(model)
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Struct {
		r.models[strings.ToLower(t.Name())] = &Model{
			name:          strings.ToLower(t.Name()),
			instance:      model,
			GetModel:      nil,
			GetModelID:    nil,
			PostModel:     nil,
			DeleteModelID: nil,
			PutModelID:    nil,
		}
	}
}

func (r *Rest) Run(addr ...string) (err error){
	path := getPath()
	port := resolveAddress(addr)
	r.Engine.GET("/", func(context *gin.Context) {
		links := gin.H{}
		for name := range r.models {
			links[name] = gin.H{
				"href" : path + port + r.BathPath + "/" + name,
			}
		}
		context.JSON(200, gin.H{
			"_links" : links,
		})
	})
	r.Engine.GET(r.BathPath + "/:model", func(context *gin.Context) {
		name := context.Param("model")
		if f := r.models[name].GetModel; f != nil {
			f(r, context)
		} else {
			ms := makeSlice(r.models[name].instance)
			r.DB.Find(ms)
			context.JSON(200, gin.H{
				"_embedded" : gin.H{
					name : ms,
				},
				"_links" : gin.H{
					"self" : gin.H{
						"href" : path + port + r.BathPath + "/" + name,
					},
				},
			})
		}
	})
	r.Engine.GET(r.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		if f := r.models[name].GetModelID; f != nil {
			f(r, context)
		} else {
			id, err := strconv.Atoi(context.Param("id"))
			m, b := r.models[name]
			if b && err == nil {
				mm := makeStruct(m.instance)
				r.DB.First(mm, id)
				context.JSON(200, mm)
			}
		}
	})
	r.Engine.POST(r.BathPath + "/:model", func(context *gin.Context) {
		name := context.Param("model")
		if f := r.models[name].PostModel; f != nil {
			f(r, context)
		} else {
			m, b := r.models[name]
			if b {
				mm := makeStruct(m.instance)
				err := context.BindJSON(mm)
				r.DB.Create(mm)
				if err == nil {
					context.JSON(200, mm)
				}
			}
		}
	})
	r.Engine.DELETE(r.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		if f := r.models[name].DeleteModelID; f != nil {
			f(r, context)
		} else {
			id, err := strconv.Atoi(context.Param("id"))
			m, b := r.models[name]
			if b && err == nil {
				mm := makeStruct(m.instance)
				r.DB.First(mm, id)
				r.DB.Delete(mm)
				context.JSON(200, gin.H{"data" : "deleted"})
			}
		}
	})
	r.Engine.PUT(r.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		if f := r.models[name].PutModelID; f != nil {
			f(r, context)
		} else {
			id, err := strconv.Atoi(context.Param("id"))
			m, b := r.models[name]
			if b && err == nil {
				mm := makeStruct(m.instance)
				r.DB.First(mm, id)
				err := context.BindJSON(mm)
				if err == nil {
					r.DB.Save(mm)
					context.JSON(200, mm)
				}
			}
		}
	})
	return r.Engine.Run(addr...)
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			return ":" + port
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

func getPath() string {
	path := "http://"
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		path = path + "localhost"
	}
	defer conn.Close()
	path = path + strings.Split(conn.LocalAddr().String(), ":")[0]
	return path
}