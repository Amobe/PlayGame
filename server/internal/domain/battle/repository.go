package battle

type Repository interface {
	Get(id string) (*Battle, error)
	Save(b *Battle) error
}
