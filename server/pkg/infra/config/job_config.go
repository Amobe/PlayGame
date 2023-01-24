package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
	"github.com/Amobe/PlayGame/server/pkg/infra/inmem"
)

const jobJsonPath = "files/job.json"

type jobConfigData struct {
	Jobs []jobData `json:"jobs"`
}

type jobData struct {
	Id                  string   `json:"id"`
	Type                string   `json:"type"`
	AvailableWeaponList []string `json:"available_weapon_list"`
}

func initJobRepository() *inmem.JobRepository {
	r := inmem.NewInmemJobRepository()

	d := jobConfigData{}
	jobConfigFile, err := configFiles.ReadFile(jobJsonPath)
	if err != nil {
		panic(fmt.Errorf("read job.json from config files: %w", err))
	}
	if err := json.Unmarshal(jobConfigFile, &d); err != nil {
		panic(fmt.Errorf("json unmarshal job config file: %w", err))
	}

	for _, data := range d.Jobs {
		job, err := jobDataToDomain(data)
		if err != nil {
			panic(fmt.Errorf("job data to domain: %w", err))
		}
		if _, err := r.Get(job.ID()); !errors.Is(err, inmem.ErrorRecordNotFound) {
			panic(fmt.Sprintf("job id is duplicated: %s", job.ID()))
		}
		_ = r.Create(job)
	}

	return r
}

func jobDataToDomain(data jobData) (vo.Job, error) {
	var availableWeapenList []vo.WeaponType
	for _, w := range data.AvailableWeaponList {
		availableWeapenList = append(availableWeapenList, vo.ToWeaponType(w))
	}
	return vo.NewJob(data.Id, data.Type, availableWeapenList)
}
