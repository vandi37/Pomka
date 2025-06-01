using UsersServiceApp;
using Grpc.Core;
using Google.Protobuf.Collections;



public class UsersHandler : Users.UsersBase
{

    private readonly IRepository repository;
    private readonly ILogger<UsersHandler> logger;

    public UsersHandler(IRepository repository, ILogger<UsersHandler> logger)
    {
        this.repository = repository;
        this.logger = logger;
    }

    public override Task<TransactionResponse> SendTransaction(TransactionRequest request, ServerCallContext context)
    {
        throw new RpcException(new Status(StatusCode.Unimplemented, ""));
    }
    public override async Task<User> Create(Common.Void request, ServerCallContext context)
    {
        try {
            var user = await repository.Create();
            return user.Grpc;
        } catch (Exception e) {
            logger.LogError(e, "Unexpected exception in Create");
            throw new RpcException(new Status(StatusCode.Internal, "Internal"));
        }
    }

    public override async Task<Common.Response> ChangeAutoBuy(Id request, ServerCallContext context)
    {
        try {
            await repository.ChangeAutoBuy(request.Id_);
            return new Common.Response();
        } catch (NotFoundException) {
            return new Common.Response{Failure = new Common.Failure{Code = Common.ErrorCode.UserNotFound}};
        }catch (Exception e) {
            logger.LogError(e, "Unexpected exception in GetUser");
            throw new RpcException(new Status(StatusCode.Internal, "Internal"));
        }
    }


    public override async Task<User> GetUser(Id request, ServerCallContext context)
    {
        try {
            var user = await repository.Get(request.Id_);
            if (user == null) throw new RpcException(new Status(StatusCode.NotFound, $"Not found transaction {request.Id_}"));
            return user.Grpc;
        } catch (Exception e) {
            logger.LogError(e, "Unexpected exception in GetUser");
            throw new RpcException(new Status(StatusCode.Internal, "Internal"));
        }
    }


    public override async Task<RepeatedUsers> GetTop(GetTopUsers request, ServerCallContext context)
    {
        try {
            var users = await repository.Top(request.Currency);
            var repeated = new RepeatedUsers();
            foreach(var u in users) {
                repeated.Users.Add(u.Grpc);
            };
            return repeated;
        } catch (InvalidCurrencyException) {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Invalid currency"));
        }catch(Exception e) {
            logger.LogError(e, "Unexpected exception in GetTop");
            throw new RpcException(new Status(StatusCode.Internal, "Internal"));
        }
    }


    public override async Task<RepeatedUsers> GetAll(Common.Void request, ServerCallContext context)
    {
        try {
            var users = await repository.All();
            var repeated = new RepeatedUsers();
            foreach(var u in users) {
                repeated.Users.Add(u.Grpc);
            };
            return repeated;
        } catch(Exception e) {
            logger.LogError(e, "Unexpected exception in GetAll");
            throw new RpcException(new Status(StatusCode.Internal, "Internal"));
        }
    }


    public override async Task<Transaction> GetTransaction(Id request, ServerCallContext context)
    {
        try {
            var transaction = await repository.GetTransaction(request.Id_);
            if (transaction == null) throw new RpcException(new Status(StatusCode.NotFound, $"Not found transaction {request.Id_}"));
            return transaction.Grpc;
        } catch (Exception e) {
            logger.LogError(e, "Unexpected exception in GetTransaction");
            throw new RpcException(new Status(StatusCode.Internal, "Internal"));
        }
    }


    public override async Task<TransactionHistory> GetTransactionHistory(Id request, ServerCallContext context)
    {
        try {
            var transactions = await repository.History(request.Id_);
            var history = new TransactionHistory();
            foreach(var t in transactions) {
                history.Transactions.Add(t.Grpc);
            };
            return history;
        } catch(Exception e) {
            logger.LogError(e, "Unexpected exception in GetTransactionHistory");
            throw new RpcException(new Status(StatusCode.Internal, "Internal"));
        }
    }


    public override async Task<TransactionHistory> GetAllTransactions(Common.Void request, ServerCallContext context)
    {
        try {
            var transactions = await repository.AllTransactions();
            var history = new TransactionHistory();
            foreach (var t in transactions) {
                history.Transactions.Add(t.Grpc);
            }
            return history;
        } catch (Exception e) {
            logger.LogError(e, "Unexpected exception in GetAllTransactions");
            throw new RpcException(new Status(StatusCode.Internal, "Internal"));
        }
    }

    public override Task<TransactionResponse> Farm(Id request, ServerCallContext context)
    {
        throw new RpcException(new Status(StatusCode.Unimplemented, ""));
    }
}
