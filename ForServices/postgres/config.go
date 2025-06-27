package postgres

type Config struct {
	Host        string
	Port        string
	User        string
	Password    string
	Database    string
	MaxAtmps    int
	DelayAtmpsS int
}
