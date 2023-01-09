package stage

type Repository interface {
	Get(id string) (*Stage, error)
}
