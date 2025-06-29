package conn

type Config struct {
	ConfigServiceUsers
	ConfigServiceCmds
}

type ConfigServiceUsers struct {
	Host string
	Port string
}

type ConfigServiceCmds struct {
	Host string
	Port string
}
