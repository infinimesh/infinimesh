package shadow

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repo interface {
	SetReported(DeviceState) (err error)
	GetReported() (DeviceState, error)
	// SetDesired(DeviceState)
	// GetDesired() (DeviceState, error)
}

type DeviceState struct {
	ID      string
	Version int64
	State   string
}

type DeviceStateDB struct {
	ID              string
	ReportedVersion int64
	ReportedState   postgres.Jsonb
	DesiredVersion  int64
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
	if err := p.db.Save(&DeviceStateDB{
		ID:              d.ID,
		ReportedVersion: d.Version,
		ReportedState:   postgres.Jsonb{[]byte(d.State)}, // nolint
	}).Error; err != nil {
		return err
	}
	return nil
}

func (p *postgresRepo) GetReported() (d DeviceState, err error) {
	if err := p.db.First(&DeviceState{}, 10).Error; err != nil {
		return DeviceState{}, err
	}
	return
}
