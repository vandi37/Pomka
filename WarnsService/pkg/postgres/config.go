package postgres

type DBConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	Database    string
	MaxAtmps    int
	DelayAtmpsS int
}
