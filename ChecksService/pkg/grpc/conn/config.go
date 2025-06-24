package conn

type Config struct {
	CfgSrvUsers ConfigServiceUsers
}

type ConfigServiceUsers struct {
	Host string
	Port string
}
