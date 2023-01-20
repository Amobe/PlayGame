package skill

type Repository interface {
	Get(id string) (Skill, error)
}
