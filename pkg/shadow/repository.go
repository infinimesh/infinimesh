package shadow

import (
	"encoding/json"

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
	if err := p.db.Table("device_state_db").Updates(&DeviceStateDB{
		ID:              d.ID,
		ReportedVersion: d.Version,
		ReportedState:   postgres.Jsonb{d.State}, // nolint
	}).Error; err != nil {
		return err
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
	if err := p.db.Table("device_state_db").Updates(&DeviceStateDB{
		ID:             d.ID,
		DesiredVersion: d.Version,
		DesiredState:   postgres.Jsonb{d.State}, // nolint
	}).Error; err != nil {
		return err
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
