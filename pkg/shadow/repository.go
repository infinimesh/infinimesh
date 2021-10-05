//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package shadow

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repo interface {
	SetReported(DeviceState) (err error)
	GetReported(id string) (DeviceState, error)
	SetDesired(DeviceState) (err error)
	GetDesired(id string) (DeviceState, error)
}

// TODO maybe use / embed DeviceStateMessage here - including Timestamp
type DeviceState struct {
	ID    string
	State DeviceStateMessage
}

type DeviceStateDB struct {
	ID              string
	ReportedVersion uint64
	ReportedState   postgres.Jsonb
	DesiredVersion  uint64
	DesiredState    postgres.Jsonb
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

func NewRedisRepo(addr string) (Repo, error) {
	return &redisRepo{
		pool: newPool(addr),
	}, nil
}

func (r *redisRepo) SetReported(d DeviceState) (err error) {
	return r.setState("reported", d)
}

func (r *redisRepo) GetReported(id string) (d DeviceState, err error) {
	return r.getState("reported", id)
}

func (r *redisRepo) getState(prefix, id string) (d DeviceState, err error) {
	conn := r.pool.Get()
	defer conn.Close()

	bytes, err := redis.Bytes(conn.Do("GET", prefix+"#"+id))
	if err != nil {
		return DeviceState{}, err
	}

	err = json.Unmarshal(bytes, &d)
	return d, err
}

func (r *redisRepo) setState(prefix string, d DeviceState) (err error) {
	conn := r.pool.Get()
	defer conn.Close()

	bytes, err := json.Marshal(&d)
	if err != nil {
		return err
	}

	err = conn.Send("SET", prefix+"#"+d.ID, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (r *redisRepo) SetDesired(d DeviceState) (err error) {
	return r.setState("desired", d)
}

func (r *redisRepo) GetDesired(id string) (d DeviceState, err error) {
	return r.getState("desired", id)
}
