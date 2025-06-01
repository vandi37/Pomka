public class Config
{
    [ConfigurationKeyName("CONN_STRING")]
    public string ConnString { get; set; } = "Host=localhost;Port=5432;Database=app;Username=user;Password=password;Ssl Mode=Disable";

    public Config(IConfiguration configurations) 
    {
        configurations.Bind(this); 
    }
    public Config() {}
} 