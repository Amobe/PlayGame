package battle

type Repository interface {
	Create(battle Battle) error
	Get(id string) (Battle, error)
	Save(stage Battle) error
	Delete(id string) error
}
