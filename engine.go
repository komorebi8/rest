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

type Engine struct {
	e        *gin.Engine
	db       *gorm.DB
	models   map[string]interface{}
	BathPath string
}

func NewEngine(e *gin.Engine, db *gorm.DB) *Engine {
	return &Engine{
		e:        e,
		db:       db,
		BathPath: "/api",
		models:   make(map[string]interface{}),  // name of model and model
	}
}

func (e *Engine) AddModel(model interface{})  {
	e.db.AutoMigrate(model)
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Struct {
		e.models[strings.ToLower(t.Name())] = model
	}
}

func (e *Engine) Run(addr ...string) (err error){
	path := getPath()
	port := resolveAddress(addr)

	e.e.GET("/", func(context *gin.Context) {
		links := gin.H{}
		for name := range e.models {
			links[name] = gin.H{
				"href" : path + port + e.BathPath + "/" + name,
			}
		}
		context.JSON(200, gin.H{
			"_links" : links,
		})
	})
	e.e.GET(e.BathPath + "/:model", func(context *gin.Context) {
		name := context.Param("model")
		ms := makeSlice(e.models[name])
		e.db.Find(ms)
		context.JSON(200, gin.H{
			"_embedded" : gin.H{
				name : ms,
			},
			"_links" : gin.H{
				"self" : gin.H{
					"href" : path + port + e.BathPath + "/" + name,
				},
			},
		})
	})
	e.e.GET(e.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		id, err := strconv.Atoi(context.Param("id"))
		m, b := e.models[name]
		if b && err == nil {
			mm := makeStruct(m)
			e.db.First(mm, id)
			context.JSON(200, mm)
		}
	})
	e.e.POST(e.BathPath + "/:model", func(context *gin.Context) {
		name := context.Param("model")
		m, b := e.models[name]
		if b {
			mm := makeStruct(m)
			err := context.BindJSON(mm)
			e.db.Create(mm)
			if err == nil {
				context.JSON(200, mm)
			}
		}
	})
	e.e.DELETE(e.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		id, err := strconv.Atoi(context.Param("id"))
		m, b := e.models[name]
		if b && err == nil {
			mm := makeStruct(m)
			e.db.First(mm, id)
			e.db.Delete(mm)
			context.JSON(200, gin.H{"data" : "deleted"})
		}
	})
	e.e.PUT(e.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		id, err := strconv.Atoi(context.Param("id"))
		m, b := e.models[name]
		if b && err == nil {
			mm := makeStruct(m)
			e.db.First(mm, id)
			err := context.BindJSON(mm)
			if err == nil {
				e.db.Save(mm)
				context.JSON(200, mm)
			}
		}
	})
	return e.e.Run(addr...)
}

// returns *[]Model
func makeSlice(model interface{}) interface{} {
	t := reflect.TypeOf(model)
	slice := reflect.MakeSlice(reflect.SliceOf(t), 10, 10)
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)
	return x.Interface()
}

// returns *Model
func makeStruct(model interface{}) interface{} {
	st := reflect.TypeOf(model)
	x := reflect.New(st)
	return x.Interface()
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
	path := ""
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		path = path + "localhost"
	}
	defer conn.Close()
	path = path + strings.Split(conn.LocalAddr().String(), ":")[0]
	return path
}