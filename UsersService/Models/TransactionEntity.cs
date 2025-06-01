

public class TransactionEntity
{
    public long Id { get; set; }
    public long? SenderId { get; set; }
    public UserEntity? Sender { get; set; }
    public long? SendAmount { get; set; }
    public Common.Currency? SendCurrency { get; set; }
    public long? ReceiverId { get; set; }
    public UserEntity? Receiver { get; set; }
    public long? ReceiveAmount { get; set; }
    public Common.Currency? ReceiveCurrency { get; set; }
    public Common.TransactionType Type { get; set; }
    public DateTime CreatedAt { get; set; } = DateTime.Now;

    public UsersServiceApp.Transaction Grpc
    {
        get {
            return new UsersServiceApp.Transaction
            {
                Id = Id,
                Sender = SenderId != null && SendAmount != null && SendCurrency != null ? new UsersServiceApp.UserTransaction
                {
                    UserId = SenderId.Value,
                    Amount = SendAmount.Value,
                    Currency = SendCurrency.Value,
                } : null,
                Receiver = ReceiverId != null && ReceiveAmount != null && ReceiveCurrency != null ? new UsersServiceApp.UserTransaction
                {
                    UserId = ReceiverId.Value,
                    Amount = ReceiveAmount.Value,
                    Currency = ReceiveCurrency.Value,
                } : null,
                Type = Type,
                CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(CreatedAt)
            };
        }
    }


}