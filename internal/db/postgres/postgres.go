package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type PGDatabase struct {
	Client *gorm.DB
}

func NewPostgresAdapter() *PGDatabase {
	return &PGDatabase{}
}

func (p *PGDatabase) Connect() error {
	// dsn (Data source name) is the connection string for the database
	// dsn := "host=" + cfg.Host + " user=" + cfg.Username + " password=" + cfg.Password + " dbname=tagore port=" + cfg.Port + " sslmode=" + cfg.SSLMode
	dsn := ""
	var err error
	p.Client, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})

	if err != nil {
		return err
	}
	return nil
}

func (p *PGDatabase) GetClient() interface{} {
	return p.Client
}

func (p *PGDatabase) Disconnect() error {
	db, err := p.Client.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
