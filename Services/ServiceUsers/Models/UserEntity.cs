
public class UserEntity
{
    public long Id {get; set;}
    public long Credits {get; set;} = 100; // Default money
    public long Stocks {get; set;} = 0;
    public UsersServiceApp.Role Role {get; set;} = UsersServiceApp.Role.Normal;
    public bool AutoBuyEnabled {get; set;} = true;
    public DateTime LastFarmingAt {get; set;}= DateTime.MinValue;
    public DateTime CreatedAt {get; set;}= DateTime.Now;

    public ICollection<TransactionEntity> SentTransactions { get; set; } = [];
    
    public ICollection<TransactionEntity> ReceivedTransactions { get; set; } = [];

    public UsersServiceApp.User Grpc {get {
        return new UsersServiceApp.User{
            Id = Id,
            Credits = Credits,
            Stocks = Stocks,
            Role = Role,
            AutoBuyEnabled = AutoBuyEnabled,
            LastFarmingAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(LastFarmingAt),
            CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(LastFarmingAt)
        };
    }}
}

