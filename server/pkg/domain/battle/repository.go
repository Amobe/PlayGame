package battle

type Repository interface {
	Get(id string) (Battle, error)
	Save(stage Battle) error
}
