package redis

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/gomodule/redigo/redis"
)

func (r *Redis) zaddHandler(content string) error {
	key := r.CurrentKey
	tmpArr := strings.Split(content, "\n")
	var args []interface{}
	temp := key
	args = append(args, key)
	for k, v := range tmpArr {
		t := strings.Split(v, SEPARATOR)
		if len(t) != 2 {
			err := errors.New("Line " + strconv.Itoa(k+1) + " include incorrect format data")
			r.Output = append(r.Output, []string{err.Error(), OUTPUT_ERROR})
			return err
		}
		temp += " " + t[1] + " " + t[0]
		args = append(args, t[1], t[0])
	}
	if err := r.delHandler(""); err != nil {
		return err
	}
	r.Output = append(r.Output, []string{"ZADD " + temp, OUTPUT_COMMAND})
	res, err := redis.Int64(r.Conn.Do("ZADD", args...))
	if err != nil {
		r.Output = append(r.Output, []string{err.Error(), OUTPUT_ERROR})
		return err
	}
	r.CurrentKey = key
	r.CurrentKeyType = TYPE_ZSET
	r.Output = append(r.Output, []string{strconv.FormatInt(res, 10), OUTPUT_RES})
	return nil
}
func (r *Redis) zscoreHandler(content string) error {
	return nil
}
func (r *Redis) zincrbyHandler(content string) error {
	return nil
}
func (r *Redis) zcardHandler(content string) error {
	r.Output = append(r.Output, []string{"ZCARD " + r.CurrentKey, OUTPUT_COMMAND})
	lenres, err := redis.Int64(r.Conn.Do("ZCARD", r.CurrentKey))
	if err != nil {
		r.Output = append(r.Output, []string{err.Error(), OUTPUT_ERROR})
		return err
	}
	r.Info = append(r.Info, []string{"zcard", strconv.FormatInt(lenres, 10)})
	return nil
}
func (r *Redis) zcountHandler(content string) error {
	return nil
}
func (r *Redis) zrangeHandler(content string) error {
	var err error
	r.Output = append(r.Output, []string{"ZRANGE " + r.CurrentKey + " 0 -1 WITHSCORES", OUTPUT_COMMAND})
	r.Detail, err = redis.StringMap(r.Conn.Do("ZRANGE", r.CurrentKey, 0, -1, "WITHSCORES"))
	if err != nil {
		r.Output = append(r.Output, []string{err.Error(), OUTPUT_ERROR})
		return err
	}
	return nil
}
func (r *Redis) zrevrangeHandler(content string) error {
	return nil
}
func (r *Redis) zrangebyscoreHandler(content string) error {
	return nil
}
func (r *Redis) zrevrangebyscoreHandler(content string) error {
	return nil
}
func (r *Redis) zrankHandler(content string) error {
	return nil
}
func (r *Redis) zrevrankHandler(content string) error {
	return nil
}
func (r *Redis) zremHandler(content string) error {
	return nil
}
func (r *Redis) zremrangebyrankHandler(content string) error {
	return nil
}
func (r *Redis) zremrangebyscoreHandler(content string) error {
	return nil
}
func (r *Redis) zrangebylexHandler(content string) error {
	return nil
}
func (r *Redis) zlexcountHandler(content string) error {
	return nil
}
func (r *Redis) zremrangebylexHandler(content string) error {
	return nil
}
func (r *Redis) zscanHandler(content string) error {
	return nil
}
func (r *Redis) zunionstoreHandler(content string) error {
	return nil
}
func (r *Redis) zinterstoreHandler(content string) error {
	return nil
}
