using Microsoft.EntityFrameworkCore;

public class UsersDbContext(DbContextOptions<UsersDbContext> options) : DbContext(options) {
    public DbSet<UserEntity> Users {get; init;}
    public DbSet<TransactionEntity> Transactions {get; init;}

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.ApplyConfiguration(new TransactionConfiguration());
        modelBuilder.ApplyConfiguration(new UserConfiguration());
        base.OnModelCreating(modelBuilder);
    }

}