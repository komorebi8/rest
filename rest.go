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
	e        *gin.Engine
	db       *gorm.DB
	models   map[string]interface{}
	BathPath string
}

func New(e *gin.Engine, db *gorm.DB) *Rest {
	return &Rest{
		e:        e,
		db:       db,
		BathPath: "/api",
		models:   make(map[string]interface{}),  // name of model and model
	}
}

func (r *Rest) AddModel(model interface{})  {
	r.db.AutoMigrate(model)
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Struct {
		r.models[strings.ToLower(t.Name())] = model
	}
}

func (r *Rest) Run(addr ...string) (err error){
	path := getPath()
	port := resolveAddress(addr)

	r.e.GET("/", func(context *gin.Context) {
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
	r.e.GET(r.BathPath + "/:model", func(context *gin.Context) {
		name := context.Param("model")
		ms := makeSlice(r.models[name])
		r.db.Find(ms)
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
	})
	r.e.GET(r.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		id, err := strconv.Atoi(context.Param("id"))
		m, b := r.models[name]
		if b && err == nil {
			mm := makeStruct(m)
			r.db.First(mm, id)
			context.JSON(200, mm)
		}
	})
	r.e.POST(r.BathPath + "/:model", func(context *gin.Context) {
		name := context.Param("model")
		m, b := r.models[name]
		if b {
			mm := makeStruct(m)
			err := context.BindJSON(mm)
			r.db.Create(mm)
			if err == nil {
				context.JSON(200, mm)
			}
		}
	})
	r.e.DELETE(r.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		id, err := strconv.Atoi(context.Param("id"))
		m, b := r.models[name]
		if b && err == nil {
			mm := makeStruct(m)
			r.db.First(mm, id)
			r.db.Delete(mm)
			context.JSON(200, gin.H{"data" : "deleted"})
		}
	})
	r.e.PUT(r.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		id, err := strconv.Atoi(context.Param("id"))
		m, b := r.models[name]
		if b && err == nil {
			mm := makeStruct(m)
			r.db.First(mm, id)
			err := context.BindJSON(mm)
			if err == nil {
				r.db.Save(mm)
				context.JSON(200, mm)
			}
		}
	})
	return r.e.Run(addr...)
}

// returns *[]Model
// Using make() to generate a slice will cause an unaddressed pointer error.
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
	path := "http://"
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		path = path + "localhost"
	}
	defer conn.Close()
	path = path + strings.Split(conn.LocalAddr().String(), ":")[0]
	return path
}