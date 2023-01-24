package config

import (
	"embed"
	_ "embed"

	"github.com/Amobe/PlayGame/server/pkg/infra/inmem"
)

//go:embed files/job.json
//go:embed files/weapon.json
var configFiles embed.FS

var (
	jobRepository    *inmem.JobRepository
	weaponRepository *inmem.WeaponRepository
)

func init() {
	jobRepository = initJobRepository()
	weaponRepository = initWeaponRepository()
}

func GetJobRepository() *inmem.JobRepository {
	return jobRepository
}

func GetWeaponRepository() *inmem.WeaponRepository {
	return weaponRepository
}
