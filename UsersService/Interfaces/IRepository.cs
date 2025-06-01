using Common;
using UsersServiceApp;

public interface IRepository {
    Task<TransactionEntity> SendTransaction(TransactionRequest transaction); 
    Task<UserEntity> Create();
    Task ChangeAutoBuy(long id);
    Task<UserEntity?> Get(long id);
    Task<List<UserEntity>> Top(Currency currency);
    Task<List<UserEntity>> All();
    Task<TransactionEntity?> GetTransaction(long id);
    Task<List<TransactionEntity>> History(long id);
    Task<List<TransactionEntity>> AllTransactions();
}   


public class InvalidCurrencyException : Exception { }
public class NotFoundException : Exception { }
public class NullTransactionException : Exception {}
public class NotEnoughMoneyException : Exception {}
public class SenderBlockedException : Exception {}
public class ReceiverBlockedException : Exception {}
public class ForbiddenException : Exception {}
public class AmountLessThanZeroException : Exception {}
public class RoleTooBigException : Exception {}