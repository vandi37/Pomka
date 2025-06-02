
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Storage.ValueConversion;
using Npgsql;

var builder = WebApplication.CreateBuilder(args);

var configurations = builder.Configuration.AddEnvironmentVariables().Build();
var config = new Config(configurations);

// Add services to the container.
builder.Services.AddDbContext<UsersDbContext>(options => {
    options.UseNpgsql(config.ConnString);
});
builder.Services.AddScoped<IRepository, UsersRepository>();
builder.Services.AddGrpc();

var app = builder.Build();

// Configure the HTTP request pipeline.
app.MapGrpcService<UsersHandler>();
app.MapGet("/", () => "Communication with gRPC endpoints must be made through a gRPC client. To learn how to create a client, visit: https://go.microsoft.com/fwlink/?linkid=2086909");

app.Run();
