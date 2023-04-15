package config

import (
	"embed"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
	"github.com/Amobe/PlayGame/server/pkg/infra/inmem"
)

//go:embed files/weapon.json
//go:embed files/skill.json
var configFiles embed.FS

const (
	weaponJsonPath = "files/weapon.json"
	skillJsonPath  = "files/skill.json"
)

type configData map[string]any

func loadConfigData(path string) ([]configData, error) {
	var data []configData
	fileData, err := configFiles.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read data from config file [%s]: %w", path, err)
	}
	if err := json.Unmarshal(fileData, &data); err != nil {
		return nil, fmt.Errorf("json unmarshal file data: %w", err)
	}
	return data, nil
}

type dataToDomainFn[T inmem.Indexer] func(data configData) (T, error)

func initRepository[T inmem.Indexer](path string, dataToDomainFn dataToDomainFn[T]) *inmem.InmemStorage[T] {
	r := inmem.NewInmemStorage[T]()
	dataList, err := loadConfigData(path)
	if err != nil {
		panic(fmt.Errorf("load config data [%s]: %w", path, err))
	}
	if err := insertDataToRepository(r, dataToDomainFn, dataList); err != nil {
		panic(fmt.Errorf("insert data to repository: %w", err))
	}
	return r
}

func insertDataToRepository[T inmem.Indexer](r *inmem.InmemStorage[T], toDomain func(data configData) (T, error), datas []configData) error {
	for _, data := range datas {
		obj, err := toDomain(data)
		if err != nil {
			return fmt.Errorf("data to domain: %w", err)
		}
		if _, err := r.Get(obj.ID()); !errors.Is(err, inmem.ErrorRecordNotFound) {
			return fmt.Errorf("id is duplicated: %s", obj.ID())
		}
		_ = r.Create(obj)
	}
	return nil
}

var (
	weaponRepository *inmem.InmemStorage[vo.Weapon]
)

func init() {
	weaponRepository = initRepository(weaponJsonPath, weaponDataToDomain)
}

func GetWeaponRepository() *inmem.InmemStorage[vo.Weapon] {
	return weaponRepository
}
