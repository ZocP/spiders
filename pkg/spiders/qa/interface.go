package qa

type Spider interface {
	Run() error
	Update() error
}
