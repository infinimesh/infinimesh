package repo

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Repo interface {
	SetDeviceStatus(DeviceState) (err error)
	GetDeviceStatus(id string) (DeviceState, error)
	DeleteDeviceStatus(id string) (err error)
}

type DeviceState struct {
	ID     string
	Status DeviceStatus
}

type redisRepo struct {
	pool *redis.Pool
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func NewRedisRepo(addr string) (repo Repo, err error) {
	repConn := &redisRepo{
		pool: newPool(addr),
	}
	return repConn, repConn.pool.Get().Err()
}
func (r *redisRepo) SetDeviceStatus(d DeviceState) (err error) {
	return r.setState(d)
}
func (r *redisRepo) GetDeviceStatus(id string) (d DeviceState, err error) {
	return r.getState(id)
}
func (r *redisRepo) DeleteDeviceStatus(id string) (err error) {
	return r.deleteState(id)
}
func (r *redisRepo) getState(devID string) (d DeviceState, err error) {
	conn := r.pool.Get()
	defer conn.Close()
	bytes, err := redis.Bytes(conn.Do("GET", devID))
	if err != nil {
		log.Printf("Error occured while getting device status from redis " + err.Error())
		return DeviceState{}, err
	}
	err = json.Unmarshal(bytes, &d)
	if err != nil {
		log.Printf("Error occured while unmarshalling the device status " + err.Error())
		return DeviceState{}, err
	}
	return d, nil
}

func (r *redisRepo) setState(d DeviceState) (err error) {
	conn := r.pool.Get()
	defer conn.Close()

	bytes, err := json.Marshal(&d)
	if err != nil {
		return err
	}

	err = conn.Send("SET", d.ID, bytes)
	if err != nil {
		log.Printf("Error occured while saving device status in redis " + err.Error())
		return err
	}
	return nil
}
func (r *redisRepo) deleteState(devID string) (err error) {
	conn := r.pool.Get()
	defer conn.Close()

	_, err = redis.Int64(conn.Do("DEL", devID))
	if err != nil {
		log.Printf("Error occured while deleting device status from redis " + err.Error())
		return err
	}
	return nil
}
