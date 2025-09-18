package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gcalvocr/go-testing/dto"
	"github.com/gcalvocr/go-testing/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBAccountRepository implements AccountRepository for MongoDB
type MongoDBAccountRepository struct {
	collection *mongo.Collection
}

// MongoDBTransactionRepository implements TransactionRepository for MongoDB
type MongoDBTransactionRepository struct {
	collection *mongo.Collection
}

// newMongoDBFactory creates MongoDB repository instances
func newMongoDBFactory(connectionString string) (*RepositoryFactory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger.Info("Connected to MongoDB database", nil)

	// Get database name from connection string or use default
	db := client.Database("bankdb")

	return &RepositoryFactory{
		AccountRepo:     &MongoDBAccountRepository{collection: db.Collection("accounts")},
		TransactionRepo: &MongoDBTransactionRepository{collection: db.Collection("transactions")},
	}, nil
}

// Account repository methods for MongoDB
func (r *MongoDBAccountRepository) Create(ctx context.Context, account *dto.AccountDTO) error {
	if account.ID == "" {
		account.ID = primitive.NewObjectID().Hex()
	}

	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, account)
	if err != nil {
		logger.Error("Failed to create account in MongoDB", err)
		return err
	}

	logger.Info("Account created in MongoDB", map[string]interface{}{
		"account_id": account.ID,
		"name":       account.Name,
	})
	return nil
}

func (r *MongoDBAccountRepository) GetByID(ctx context.Context, id string) (*dto.AccountDTO, error) {
	var account dto.AccountDTO
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Account not found
		}
		logger.Error("Failed to get account from MongoDB", err)
		return nil, err
	}
	return &account, nil
}

func (r *MongoDBAccountRepository) GetAll(ctx context.Context) ([]*dto.AccountDTO, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to query accounts from MongoDB", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []*dto.AccountDTO
	for cursor.Next(ctx) {
		var account dto.AccountDTO
		if err := cursor.Decode(&account); err != nil {
			logger.Error("Failed to decode account from MongoDB", err)
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	if err := cursor.Err(); err != nil {
		logger.Error("Cursor error in MongoDB", err)
		return nil, err
	}

	return accounts, nil
}

func (r *MongoDBAccountRepository) Update(ctx context.Context, id string, update *dto.UpdateAccountRequest) error {
	updateDoc := bson.M{"updated_at": time.Now()}

	if update.Name != nil {
		updateDoc["name"] = *update.Name
	}
	if update.Balance != nil {
		updateDoc["balance"] = *update.Balance
	}
	if update.Currency != nil {
		updateDoc["currency"] = *update.Currency
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateDoc})
	if err != nil {
		logger.Error("Failed to update account in MongoDB", err)
		return err
	}

	logger.Info("Account updated in MongoDB", map[string]interface{}{
		"account_id": id,
	})
	return nil
}

func (r *MongoDBAccountRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		logger.Error("Failed to delete account from MongoDB", err)
		return err
	}

	logger.Info("Account deleted from MongoDB", map[string]interface{}{
		"account_id":    id,
		"deleted_count": result.DeletedCount,
	})
	return nil
}

func (r *MongoDBAccountRepository) GetByName(ctx context.Context, name string) (*dto.AccountDTO, error) {
	var account dto.AccountDTO
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Account not found
		}
		logger.Error("Failed to get account by name from MongoDB", err)
		return nil, err
	}
	return &account, nil
}

func (r *MongoDBAccountRepository) UpdateBalance(ctx context.Context, id string, newBalance float64) error {
	updateDoc := bson.M{
		"balance":    newBalance,
		"updated_at": time.Now(),
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateDoc})
	if err != nil {
		logger.Error("Failed to update account balance in MongoDB", err)
		return err
	}

	logger.Info("Account balance updated in MongoDB", map[string]interface{}{
		"account_id":  id,
		"new_balance": newBalance,
	})
	return nil
}

// Transaction repository methods for MongoDB
func (r *MongoDBTransactionRepository) Create(ctx context.Context, transaction *dto.TransactionDTO) error {
	if transaction.ID == "" {
		transaction.ID = primitive.NewObjectID().Hex()
	}

	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, transaction)
	if err != nil {
		logger.Error("Failed to create transaction in MongoDB", err)
		return err
	}

	logger.Info("Transaction created in MongoDB", map[string]interface{}{
		"transaction_id": transaction.ID,
		"account_id":     transaction.AccountID,
		"amount":         transaction.Amount,
		"type":           transaction.Type,
	})
	return nil
}

func (r *MongoDBTransactionRepository) GetByID(ctx context.Context, id string) (*dto.TransactionDTO, error) {
	var transaction dto.TransactionDTO
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Transaction not found
		}
		logger.Error("Failed to get transaction from MongoDB", err)
		return nil, err
	}
	return &transaction, nil
}

func (r *MongoDBTransactionRepository) GetByAccountID(ctx context.Context, accountID string) ([]*dto.TransactionDTO, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"account_id": accountID})
	if err != nil {
		logger.Error("Failed to query transactions by account ID from MongoDB", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*dto.TransactionDTO
	for cursor.Next(ctx) {
		var transaction dto.TransactionDTO
		if err := cursor.Decode(&transaction); err != nil {
			logger.Error("Failed to decode transaction from MongoDB", err)
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	if err := cursor.Err(); err != nil {
		logger.Error("Cursor error in MongoDB", err)
		return nil, err
	}

	return transactions, nil
}

func (r *MongoDBTransactionRepository) GetAll(ctx context.Context) ([]*dto.TransactionDTO, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to query all transactions from MongoDB", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []*dto.TransactionDTO
	for cursor.Next(ctx) {
		var transaction dto.TransactionDTO
		if err := cursor.Decode(&transaction); err != nil {
			logger.Error("Failed to decode transaction from MongoDB", err)
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	if err := cursor.Err(); err != nil {
		logger.Error("Cursor error in MongoDB", err)
		return nil, err
	}

	return transactions, nil
}

func (r *MongoDBTransactionRepository) Update(ctx context.Context, id string, transaction *dto.TransactionDTO) error {
	transaction.UpdatedAt = time.Now()

	updateDoc := bson.M{
		"account_id": transaction.AccountID,
		"amount":     transaction.Amount,
		"type":       transaction.Type,
		"updated_at": transaction.UpdatedAt,
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateDoc})
	if err != nil {
		logger.Error("Failed to update transaction in MongoDB", err)
		return err
	}

	logger.Info("Transaction updated in MongoDB", map[string]interface{}{
		"transaction_id": id,
	})
	return nil
}

func (r *MongoDBTransactionRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		logger.Error("Failed to delete transaction from MongoDB", err)
		return err
	}

	logger.Info("Transaction deleted from MongoDB", map[string]interface{}{
		"transaction_id": id,
		"deleted_count":  result.DeletedCount,
	})
	return nil
}

func (r *MongoDBTransactionRepository) GetTransactionSummary(ctx context.Context, accountID string) (*dto.TransactionSummary, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"account_id": accountID}}},
		{{"$group", bson.M{
			"_id":                "$account_id",
			"total_transactions": bson.M{"$sum": 1},
			"total_deposits": bson.M{"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []string{"$type", "deposit"}},
					"$amount",
					0,
				},
			}},
			"total_withdrawals": bson.M{"$sum": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []string{"$type", "withdrawal"}},
					"$amount",
					0,
				},
			}},
			"last_transaction_at": bson.M{"$max": "$created_at"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Error("Failed to aggregate transaction summary from MongoDB", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var summary dto.TransactionSummary
	summary.AccountID = accountID

	if cursor.Next(ctx) {
		var result struct {
			TotalTransactions int       `bson:"total_transactions"`
			TotalDeposits     float64   `bson:"total_deposits"`
			TotalWithdrawals  float64   `bson:"total_withdrawals"`
			LastTransactionAt time.Time `bson:"last_transaction_at"`
		}

		if err := cursor.Decode(&result); err != nil {
			logger.Error("Failed to decode transaction summary from MongoDB", err)
			return nil, err
		}

		summary.TotalTransactions = result.TotalTransactions
		summary.TotalDeposits = result.TotalDeposits
		summary.TotalWithdrawals = result.TotalWithdrawals
		summary.LastTransactionAt = &result.LastTransactionAt
	}

	// Get current balance from accounts collection
	accountCollection := r.collection.Database().Collection("accounts")
	var account struct {
		Balance float64 `bson:"balance"`
	}
	err = accountCollection.FindOne(ctx, bson.M{"_id": accountID}).Decode(&account)
	if err != nil {
		logger.Error("Failed to get current balance from MongoDB", err)
		return nil, err
	}

	summary.CurrentBalance = account.Balance
	return &summary, nil
}
