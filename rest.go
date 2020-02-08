package rest

import (
	"errors"
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
		// todo paging and sorting
		name := context.Param("model")
		if f := r.models[name].GetModelFunc; f != nil {
			f(r, context)
		} else {
			r.models[name].OperateInstanceSlice(func(ms interface{}) {
				if err := r.DB.Find(ms).Error; err == nil {
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
		}
	})
	r.Engine.GET(r.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		if f := r.models[name].GetModelIDFunc; f != nil {
			f(r, context)
		} else {
			r.models[name].OperateInstance(func(mm interface{}) {
				if id, err := strconv.Atoi(context.Param("id")); err == nil {
					if err = r.DB.First(mm, id).Error; err == nil {
						context.JSON(200, mm)
					}
				}
			})
		}
	})
	r.Engine.POST(r.BathPath + "/:model", func(context *gin.Context) {
		name := context.Param("model")
		if f := r.models[name].PostModelFunc; f != nil {
			f(r, context)
		} else {
			r.models[name].OperateInstance(func(mm interface{}) {
				if err := context.BindJSON(mm); err == nil {
					if err = r.DB.Create(mm).Error; err == nil {
						context.JSON(200, mm)
					}
				}
			})
		}
	})
	r.Engine.DELETE(r.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		if f := r.models[name].DeleteModelIDFunc; f != nil {
			f(r, context)
		} else {
			r.models[name].OperateInstance(func(mm interface{}) {
				if id, err := strconv.Atoi(context.Param("id")); err == nil {
					if err = r.DB.First(mm, id).Error; err == nil {
						if err = r.DB.Delete(mm).Error; err == nil {
							context.JSON(200, gin.H{
								"data" : "deleted",
							})
						}
					}
				}
			})
		}
	})
	r.Engine.PUT(r.BathPath + "/:model/:id", func(context *gin.Context) {
		name := context.Param("model")
		if f := r.models[name].PutModelIDFunc; f != nil {
			f(r, context)
		} else {
			r.models[name].OperateInstance(func(mm interface{}) {
				if id, err := strconv.Atoi(context.Param("id")); err == nil {
					if err = r.DB.First(mm, id).Error; err == nil {
						if err = context.BindJSON(mm); err == nil {
							if err = r.DB.Save(mm).Error; err == nil {
								context.JSON(200, mm)
							}
						}
					}
				}
			})
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