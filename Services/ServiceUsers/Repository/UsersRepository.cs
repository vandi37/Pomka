using Common;
using Microsoft.EntityFrameworkCore;
using UsersServiceApp;

public class UsersRepository : IRepository
{
    public const int top = 10;
    public UsersRepository(UsersDbContext db)
    {
        this.db = db;
    }
    private readonly UsersDbContext db;

    public async Task<TransactionEntity> SendTransaction(TransactionRequest transaction)
    {
        UserTransaction? sender = transaction.Sender;
        UserTransaction? receiver = transaction.Receiver;
        if (sender == null && receiver == null) throw new NullTransactionException();
        if (sender != null && await db.Users.AnyAsync(u => u.Id == receiver.UserId && u.Role == Role.Blocked)) throw new SenderBlockedException();
        if (receiver != null && await db.Users.AnyAsync(u => u.Id == receiver.UserId && u.Role == Role.Blocked)) throw new ReceiverBlockedException();
        var t = new TransactionEntity { Type = transaction.Type };
        
        if (sender != null) {
            t.SenderId = sender.UserId;
            t.SendAmount = sender.Amount;
            t.SendCurrency = sender.Currency;
            var senderProfile = await db.Users.FirstOrDefaultAsync(u => u.Id == sender.UserId) ?? throw new NotFoundException();
            switch (transaction.Type) {
                case TransactionType.Block | TransactionType.User | TransactionType.Moderator | TransactionType.Warn:
                    if (senderProfile.Role < Role.Moderator) throw new ForbiddenException();
                    break;
                case TransactionType.Set | TransactionType.Get | TransactionType.CreatePromoCode:
                    if (senderProfile.Role != Role.Creator) throw new ForbiddenException();
                    break;
            }
            if (sender.Currency == Currency.Credits) senderProfile.Credits -= sender.Amount;
            else if (sender.Currency == Currency.Stocks) senderProfile.Stocks -= sender.Amount;
        }
        if (receiver != null) {
            t.ReceiverId = receiver.UserId;
            t.ReceiveAmount = receiver.Amount;
            t.ReceiveCurrency = receiver.Currency;
            var receiverProfile = await db.Users.FirstOrDefaultAsync(u => u.Id == receiver.UserId) ?? throw new NotFoundException();
            switch (transaction.Type) {
                
            }
            
        }
    }

    public async Task<UserEntity> Create()
    {
        var user = new UserEntity();
        await db.Users.AddAsync(user);
        await db.SaveChangesAsync();
        return user;
    }

    public async Task ChangeAutoBuy(long id)
    {
        var rowsAffected = await db.Users.Where(u => u.Id == id).ExecuteUpdateAsync(u => u.SetProperty(x => x.AutoBuyEnabled, x => !x.AutoBuyEnabled));
        if (rowsAffected <= 0) throw new NotFoundException();
    }

    public async Task<List<UserEntity>> Top(Currency currency)
    {
        var query = db.Users.AsNoTracking();
        switch (currency)
        {
            case Currency.Credits:
                query = query.OrderByDescending(u => u.Credits);
                break;
            case Currency.Stocks:
                query = query.OrderByDescending(u => u.Stocks);
                break;
            default:
                throw new InvalidCurrencyException();
        }
        return await query.Take(top).ToListAsync();

    }

    public async Task<List<UserEntity>> All()
    {
        return await db.Users.AsNoTracking().ToListAsync();
    }

    public async Task<TransactionEntity?> GetTransaction(long id)
    {
        return await db.Transactions.AsNoTracking().FirstOrDefaultAsync(t => t.Id == id);
    }

    public async Task<List<TransactionEntity>> History(long id)
    {
        return await db.Transactions.AsNoTracking().Where(t => t.SenderId == id || t.ReceiverId == id).ToListAsync();
    }

    public async Task<List<TransactionEntity>> AllTransactions()
    {
        return await db.Transactions.AsNoTracking().ToListAsync();
    }

    public async Task<UserEntity?> Get(long id)
    {
        return await db.Users.FirstOrDefaultAsync(u => u.Id == id);
    }
}

