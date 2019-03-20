package shadow

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repo interface {
	SetReported(DeviceState) (err error)
	GetReported(id string) (DeviceState, error)
	SetDesired(DeviceState) (err error)
	GetDesired(id string) (DeviceState, error)
}

type DeviceState struct {
	ID      string
	Version uint64
	State   json.RawMessage
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

type postgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(addr string) (Repo, error) {
	db, err := gorm.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.LogMode(false)
	db.SingularTable(true)
	db.AutoMigrate(&DeviceStateDB{})

	return &postgresRepo{
		db: db,
	}, nil
}

func (p *postgresRepo) SetReported(d DeviceState) (err error) {
	update := DeviceStateDB{
		ID:              d.ID,
		ReportedVersion: d.Version,
		ReportedState:   postgres.Jsonb{d.State}, // nolint
	}
	if result := p.db.Model(&update).Updates(update); result.Error != nil {
		return err
	} else {
		if result.RowsAffected == 0 {
			if err := p.db.Create(&update).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *postgresRepo) GetReported(id string) (DeviceState, error) {
	var d DeviceStateDB
	if err := p.db.First(&d, "id = ?", id).Error; err != nil {
		return DeviceState{}, err
	}
	return DeviceState{ID: d.ID,
		Version: d.ReportedVersion,
		State:   d.ReportedState.RawMessage,
	}, nil
}

func (p *postgresRepo) SetDesired(d DeviceState) (err error) {
	update := DeviceStateDB{
		ID:             d.ID,
		DesiredVersion: d.Version,
		DesiredState:   postgres.Jsonb{d.State}, // nolint
	}
	if result := p.db.Model(&update).Updates(update); result.Error != nil {
		return err
	} else {
		if result.RowsAffected == 0 {
			if err := p.db.Create(&update).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *postgresRepo) GetDesired(id string) (DeviceState, error) {
	var d DeviceStateDB
	if err := p.db.First(&d, "id = ?", id).Error; err != nil {
		return DeviceState{}, err
	}
	return DeviceState{ID: d.ID,
		Version: d.DesiredVersion,
		State:   d.DesiredState.RawMessage,
	}, nil
}
