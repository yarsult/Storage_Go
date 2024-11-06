package server

import (
	"context"
	"fmt"
	"net/http"
	"proj1/internal/pkg/storage"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host    string
	storage *storage.SliceStorage
	engine  *gin.Engine
	server  *http.Server
}

type Entry struct {
	Value string `json:"value"`
}

func New(host string, st *storage.SliceStorage) *Server {
	engine := gin.New()
	s := &Server{
		host:    host,
		storage: st,
		engine:  engine,
		server: &http.Server{
			Addr:    host,
			Handler: engine,
		},
	}
	s.registerRoutes()
	return s
}

func (r *Server) registerRoutes() {
	r.engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	scalar := r.engine.Group("/scalar")
	{
		scalar.POST("set/:key/:value", r.handlerSet)
		scalar.GET("get/:key", r.handlerGet)
	}

	mapg := r.engine.Group("/map")
	{
		mapg.POST("hset/:key", r.handlerHSet)
		mapg.GET("hget/:key/:field", r.handlerHGet)
	}

	slice := r.engine.Group("/slice")
	{
		slice.POST("lpush/:key", r.handlerLPush)
		slice.POST("rpush/:key", r.handlerRPush)
		slice.POST("raddtoset/:key", r.handlerRAddToSet)
		slice.POST("/slice/lset/:key/:index/:elem", r.handlerLSet)
		slice.GET("lpop/:key", r.handlerLPop)
		slice.GET("rpop/:key", r.handlerRPop)
		slice.GET("/slice/lget/:key/:index", r.handlerLGet)
	}
	r.engine.POST("/any/expire/:key/:seconds", r.handlerExpire)
	r.engine.GET("/keys/:exp", r.handlerRegExpKeys)
}

func (r *Server) handlerSet(ctx *gin.Context) {
	key := ctx.Param("key")
	value := ctx.Param("value")
	err := r.storage.Set(key, value)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set key"})
		return
	}

	exp := ctx.Query("exp")
	if exp != "" {
		tmp, err := strconv.ParseInt(exp, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "uncorrect expiration time"})
			return
		}

		r.storage.Expire(key, tmp)
	}

	ctx.Status(http.StatusOK)
}

func (r *Server) handlerGet(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		return
	}

	v, ok := r.storage.Get(key)
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
		return
	}

	field := ctx.Param("field")
	if field == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := r.storage.HGet(key, field)
	if err != nil || res == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value: *res})
}

func (r *Server) handlerHSet(ctx *gin.Context) {
	key := ctx.Param("key")
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
}

func (r *Server) handlerLPush(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		return
	}

	var vals []string
	if err := ctx.Bind(&vals); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.LPush(key, vals)
	ctx.Status(http.StatusOK)
}

func (r *Server) handlerRPush(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		return
	}

	var vals []string
	if err := ctx.Bind(&vals); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	r.storage.RPush(key, vals)
	ctx.Status(http.StatusOK)
}

func (r *Server) handlerRAddToSet(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
		return
	}

	var vals []string
	if err := ctx.Bind(&vals); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r.storage.RAddToSet(key, vals)
	ctx.Status(http.StatusOK)
}

func (r *Server) handlerLPop(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
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
}

func (r *Server) handlerRPop(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
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
}

func (r *Server) handlerLSet(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
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
}

func (r *Server) handlerLGet(ctx *gin.Context) {
	key := ctx.Param("key")
	if r.storage.CheckIfExpired(key) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "element has expired"})
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
}

func (r *Server) handlerExpire(ctx *gin.Context) {
	key := ctx.Param("key")
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
}

func (r *Server) handlerRegExpKeys(ctx *gin.Context) {
	exp := ctx.Param("exp")
	res, err := r.storage.RegExKeys(exp)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid expression"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (r *Server) Start() error {
	fmt.Println("Starting server at", r.host)
	if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}

func (r *Server) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down server...")

	if err := r.storage.SaveToFile(r.storage.Path); err != nil {
		fmt.Println("Error saving storage:", err)
		return err
	}

	if err := r.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server forced to shutdown: %w", err)
	}

	fmt.Println("Server exited gracefully")
	return nil
}
