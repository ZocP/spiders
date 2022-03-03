package config

//spider storage location

const (
	LocalTxt = 0
	LocalCSV = 1
	RemoteDB = 2
)

//service id
type ServiceID int

const (
	QueryQA ServiceID = 1
)

type DB struct {
	URL      string `json:"url"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
	DBName   string `json:"db_name"`
}

func (c *Config) GetServiceDB(service ServiceID) DB {
	switch service {
	case QueryQA:
		return c.Services.QueryQA.DB
		// can add more services here to get DB
	}
	return DB{}
}
