package nsfContext

import (
	"nsf/internal/db"
	"nsf/logger"
)

/*
Important Note:
- Do not modify this file.
- All you need to do is to implement the customDb in the internal/db package
*/

type dbContextIE struct {
	dbType string

	*logger.NsfLogger
}

type dbContext struct {
	db db.DbIf

	*logger.NsfLogger
}

func newDbContext(dbContextIE *dbContextIE) (*dbContext, error) {
	db, err := db.NewDb(dbContextIE.dbType)
	if err != nil {
		return nil, err
	}

	return &dbContext{
		db: db,

		NsfLogger: dbContextIE.NsfLogger,
	}, nil
}

func (d *dbContext) release() {
	d.NsfLogger.DbLog.Infoln("Release dbContext...")

	if err := d.db.Release(); err != nil {
		d.NsfLogger.DbLog.Errorf("Failed to release dbContext: %v", err)
	}

	d.NsfLogger.DbLog.Infoln("dbContext released")
}

func (d *dbContext) Exist(key string) (bool, error) {
	return d.db.Exist(key)
}

func (d *dbContext) Save(key, value string) error {
	return d.db.Save(key, value)
}

func (d *dbContext) SaveAll(data map[string]string) error {
	return d.db.SaveAll(data)
}

func (d *dbContext) Load(key string) (string, error) {
	return d.db.Load(key)
}

func (d *dbContext) LoadAll() (map[string]string, error) {
	return d.db.LoadAll()
}

func (d *dbContext) Remove(key string) error {
	return d.db.Remove(key)
}

func (d *dbContext) RemoveAll() error {
	return d.db.RemoveAll()
}
