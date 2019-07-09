package controller

import (
	"fmt"
	"github.com/inosy22/golang-echo-try/myredis"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Request struct {
	Key   string `form:"key"`
	Value string `form:"value"`
	Ttl   int    `form:"ttl"`
}

type Response struct {
	Command string `json:"command"`
	Key     string `json:"key"`
	Result  string `json:"result"`
}

func RedisHello(c echo.Context) error {
	_, err := myredis.NewRedisClient()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failure, Connect Redis!")
	}

	// HTML File Open (Templateではないので無理矢理文字列化する)
	f, err := os.Open("../views/redis.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failure, Open HTML file!")
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failure, Read HTML file!")
	}

	return c.HTML(http.StatusOK, string(b))
}

func RedisGet(c echo.Context) error {
	key := c.Param("key")
	redis, err := myredis.NewRedisClient()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failure, Connect Redis!")
	}
	cmd := redis.Get(key)
	val, err := cmd.Result()
	json := &Response{
		Command: fmt.Sprintf("Get %s", key),
		Key:     key,
		Result:  val,
	}
	return c.JSON(http.StatusOK, json)
}

func RedisPut(c echo.Context) error {
	key := c.Param("key")
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Failure, Request parameter is invalid")
	}
	req.Key = key
	return set(c, req)
}

func RedisPost(c echo.Context) error {
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Failure, Request parameter is invalid")
	}

	return set(c, req)
}

func RedisDelete(c echo.Context) error {
	key := c.Param("key")

	redis, err := myredis.NewRedisClient()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failure, Connect Redis!")
	}

	cmd := redis.Del(key)
	res, err := cmd.Result()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failure, Del %s to Redis!", key))
	}
	json := &Response{
		Command: fmt.Sprintf("Del %s", key),
		Key:     key,
		Result:  strconv.Itoa(int(res)),
	}
	return c.JSON(http.StatusOK, json)
}

func set(c echo.Context, req *Request) error {
	redis, err := myredis.NewRedisClient()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failure, Connect Redis!")
	}
	d, err := time.ParseDuration(strconv.Itoa(req.Ttl))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failure, Ttl is invalid!")
	}
	cmd := redis.Set(req.Key, req.Value, d)
	res, err := cmd.Result()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failure, Set %s:%s to Redis!", req.Key, req.Value))
	}
	json := &Response{
		Command: fmt.Sprintf("Set %s %s", req.Key, req.Value),
		Key:     req.Key,
		Result:  res,
	}
	return c.JSON(http.StatusOK, json)
}
