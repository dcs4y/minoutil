package system

import (
	"fmt"
	"game/entity"
	"game/utils/dbclient"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	DB = dbclient.DB
	DB.AutoMigrate(&entity.Job{})
}

func GetJobList(param entity.Job) (jobs []entity.Job) {
	dbclient.DB.Find(&jobs, param)
	return jobs
}

func GetJobById(id uint64) *entity.Job {
	o, b := dbclient.DB.Model(&entity.Job{}).Get(fmt.Sprintf("%d", id))
	if b {
		return o.(*entity.Job)
	}
	return nil
}

func SaveJobById(job *entity.Job) error {
	return dbclient.DB.Save(job).Error
}
