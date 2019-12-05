package model

type ServiceConfig struct {
	Product Onlinestruct `yaml:"Product"`
	Boe     Onlinestruct `yaml:"Boe"`
	Local   Localstruct  `yaml:"Local"`
}

type Onlinestruct struct {
	Mysql MysqlConfigOnline `yaml:"Databaseconfig"`
}

type Localstruct struct {
	Mysql MysqlConfigLocal `yaml:"Databaseconfig"`
}

type MysqlConfigLocal struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type MysqlConfigOnline struct {
	ConsulParm string `yaml:"consulparm"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Database   string `yaml:"database"`
}
