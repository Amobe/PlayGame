package character

type Repository interface {
	Create(character Character) error
	Get(id string) (Character, error)
	Save(character Character) error
	Delete(id string) error
}
