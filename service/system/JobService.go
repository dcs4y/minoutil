package system

import (
	"fmt"
	"game/entity/model"
	"game/utils/dbclient"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	DB = dbclient.DB
	DB.AutoMigrate(&model.Job{})
}

func GetJobList(param model.Job) (jobs []model.Job) {
	dbclient.DB.Find(&jobs, param)
	return jobs
}

func GetJobById(id uint64) *model.Job {
	o, b := dbclient.DB.Model(&model.Job{}).Get(fmt.Sprintf("%d", id))
	if b {
		return o.(*model.Job)
	}
	return nil
}

func SaveJobById(job *model.Job) error {
	return dbclient.DB.Save(job).Error
}
