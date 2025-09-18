# Bank API Postman Collection

This directory contains Postman collections and environments for testing the Bank API.

## Files

- `Bank_API_Collection.postman_collection.json` - Complete collection with all API endpoints
- `Bank_API_Environment.postman_environment.json` - Environment configuration for local development

## How to Import

1. Open Postman
2. Click "Import" in the top left
3. Select "File"
4. Import both files:
   - First import the environment file
   - Then import the collection file

## Available Endpoints

### Health Check
- **GET** `/health` - Check API health status

### Accounts
- **GET** `/accounts` - Get all accounts
- **POST** `/accounts` - Create a new account
- **GET** `/accounts/{id}` - Get account by ID

### Transactions
- **GET** `/accounts/{account_id}/transactions` - Get transactions for an account
- **POST** `/transactions` - Create a new transaction (deposit/withdrawal)

### Exchange Rates
- **GET** `/exchange?from=USD&to=EUR` - Get exchange rate between currencies

## Environment Variables

The collection uses the following environment variable:
- `baseUrl` - Set to `http://localhost:8080` for local development

## Testing Workflow

1. Start the application: `make up`
2. Import the Postman collection and environment
3. Test the health endpoint first
4. Create an account using the POST `/accounts` endpoint
5. Create transactions using the POST `/transactions` endpoint
6. Test exchange rates with the GET `/exchange` endpoint

## Sample Requests

### Create Account
```json
{
  "name": "John Doe",
  "balance": 1000.00,
  "currency": "USD"
}
```

### Create Transaction
```json
{
  "account_id": 1,
  "amount": 500.00,
  "type": "deposit"
}
```

### Get Exchange Rate
```
GET /exchange?from=USD&to=EUR
```

## Response Examples

All endpoints return JSON responses. Check the collection for detailed examples of successful and error responses.

## Notes

- The API uses PostgreSQL for data persistence
- Transactions are atomic and update account balances
- Exchange rates are fetched from an external API
- All requests are logged with structured logging
- Error responses include appropriate HTTP status codes

## Troubleshooting

- Make sure the application is running (`make up`)
- Check the logs with `make logs`
- Verify the base URL in the environment is correct
- Ensure PostgreSQL is running and accessible