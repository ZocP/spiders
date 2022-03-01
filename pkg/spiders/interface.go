package spiders


type Spider interface{
	Run() error
	Update() error
}