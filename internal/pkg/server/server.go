package server

import (
	"fmt"
	"net/http"
	"proj1/internal/pkg/storage"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host    string
	storage *storage.SliceStorage
}

type Entry struct {
	Value string `json:"value"`
}

func New(host string, st *storage.SliceStorage) *Server {
	s := &Server{
		host:    host,
		storage: st,
	}

	return s
}

func (r *Server) newAPI() *gin.Engine {
	engine := gin.New()

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	engine.POST("/scalar/set/:key/:value", r.handlerSet)
	engine.GET("/scalar/get/:key", r.handlerGet)

	engine.POST("/map/hset/:key", r.handlerHSet)
	engine.GET("/map/hget/:key/:field", r.handlerHGet)

	engine.POST("/slice/lpush/:key", r.handlerLPush)
	engine.POST("/slice/rpush/:key", r.handlerRPush)
	engine.POST("/slice/raddtoset/:key", r.handlerRAddToSet)
	engine.POST("/slice/lset/:key/:index/:elem", r.handlerLSet)
	engine.GET("/slice/lpop/:key", r.handlerLPop)
	engine.GET("/slice/rpop/:key", r.handlerRPop)
	engine.GET("/slice/lget/:key/:index", r.handlerLGet)

	engine.POST("/any/expire/:key/:seconds", r.handlerExpire)
	return engine
}

func (r *Server) handlerSet(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	value := ctx.Param("value")
	err := r.storage.Set(key, value)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	exp := ctx.Query("exp")
	if exp != "" {
		tmp, err := strconv.ParseInt(exp, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "uncorrect expiration time"})
			r.storage.SaveToFile("slice_storage.json")
			return
		}
		r.storage.Expire(key, tmp)
	}
	ctx.Status(http.StatusOK)
	r.storage.SaveToFile("slice_storage.json")
}

func (r *Server) handlerGet(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	v, ok := r.storage.Get(key)
	fmt.Println(ok)
	if !ok {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value: v})
}

func (r *Server) handlerHGet(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	field := ctx.Param("field")
	res, err := r.storage.HGet(key, field)
	if err != nil || res == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, Entry{Value: *res})
}

func (r *Server) handlerHSet(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	var maps []map[string]string
	if err := ctx.Bind(&maps); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c, err := r.storage.HSet(key, maps)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, c)
	r.storage.SaveToFile("slice_storage.json")
}

func (r *Server) handlerLPush(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	var vals []string
	if err := ctx.Bind(&vals); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	r.storage.LPush(key, vals)
	ctx.Status(http.StatusOK)
	r.storage.SaveToFile("slice_storage.json")
}

func (r *Server) handlerRPush(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	var vals []string
	if err := ctx.Bind(&vals); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	r.storage.RPush(key, vals)
	ctx.Status(http.StatusOK)
	r.storage.SaveToFile("slice_storage.json")
}

func (r *Server) handlerRAddToSet(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	var vals []string
	if err := ctx.Bind(&vals); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	r.storage.RAddToSet(key, vals)
	ctx.Status(http.StatusOK)
	r.storage.SaveToFile("slice_storage.json")
}

func (r *Server) handlerLPop(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	startstr := ctx.Query("start")
	if startstr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start index is required"})
		return
	}
	start, err := strconv.Atoi(startstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start index"})
		return
	}
	endstr := ctx.Query("end")
	var indexes []int
	indexes = append(indexes, start)
	if endstr != "" {
		end, err := strconv.Atoi(endstr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end index"})
			return
		}
		indexes = append(indexes, end)
	}
	result := r.storage.LPop(key, indexes...)

	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no elements found or uncorrect indexes"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"result": result})
	r.storage.SaveToFile("slice_storage.json")
}

func (r *Server) handlerRPop(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	startstr := ctx.Query("start")
	if startstr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start index is required"})
		return
	}
	start, err := strconv.Atoi(startstr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start index"})
		return
	}
	endstr := ctx.Query("end")
	var indexes []int
	indexes = append(indexes, start)
	if endstr != "" {
		end, err := strconv.Atoi(endstr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end index"})
			return
		}
		indexes = append(indexes, end)
	}
	result := r.storage.LPop(key, indexes...)

	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no elements found or uncorrect indexes"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"result": result})
	r.storage.SaveToFile("slice_storage.json")
}

func (r *Server) handlerLSet(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	ind, err := strconv.Atoi(ctx.Param("index"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "index must be integer"})
	}
	elem := ctx.Param("elem")
	_, err = r.storage.LSet(key, ind, elem)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid index"})
		return
	}
	ctx.Status(http.StatusOK)
	r.storage.SaveToFile("slice_storage.json")

}

func (r *Server) handlerLGet(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	ind, err := strconv.Atoi(ctx.Param("index"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "index must be integer"})
	}
	res, err := r.storage.LGet(key, ind)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid index"})
		return
	}
	ctx.JSON(http.StatusOK, res)
	r.storage.SaveToFile("slice_storage.json")
}

func (r *Server) handlerExpire(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		r.storage.SaveToFile("slice_storage.json")
		return
	}
	seconds, err := strconv.ParseInt(ctx.Param("seconds"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid time"})
		return
	}
	res := r.storage.Expire(key, seconds)
	if res == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid key"})
		return
	}
	ctx.JSON(http.StatusOK, res)
	r.storage.SaveToFile("slice_storage.json")
}

func (r *Server) Start() {
	r.newAPI().Run(r.host)
}
