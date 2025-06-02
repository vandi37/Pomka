using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

public class TransactionConfiguration : IEntityTypeConfiguration<TransactionEntity>
{
    public void Configure(EntityTypeBuilder<TransactionEntity> builder)
    {
        builder.HasKey(t => t.Id);
        builder
            .HasOne(t => t.Sender)
            .WithMany(u => u. SentTransactions)
            .HasForeignKey(t => t.SenderId)
            .IsRequired(false)
            .OnDelete(DeleteBehavior.NoAction);
        builder
            .HasOne(t => t.Receiver)
            .WithMany(u => u. ReceivedTransactions)
            .HasForeignKey(t => t.ReceiverId)
            .IsRequired(false)
            .OnDelete(DeleteBehavior.NoAction);
    }
}