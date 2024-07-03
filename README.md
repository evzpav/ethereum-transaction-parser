# Ethereum Transactions Parser

It is a server written in Golang, which run in intervals of 10 seconds(variable) to get all transactions from latest parsed block until the lastest block in Ethereum blockchain.

## Running
```bash
    go run main.go

```

It has a REST API running on http://localhost:8787 with 3 endpoints:

1. POST /subscribe

    Request body format:
    ```
        {
           "address": "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5"
        }
    ```  
    Response: 
    ```
        {
	        "message": "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5 subscribed"
        }
    ```

    Example:
    ```
    curl --request POST \
    --url http://localhost:8787/subscribe \
    --header 'Content-Type: application/json' \
    --data '{
    "address": "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5"
    }'
    ```

 2. GET /transactions with query param `address`

    Example: 
    ```
        curl --request GET \
        --url 'http://localhost:8787/transactions?address=0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5'

    ```

    Response:
    ```
    {
	    "transactions": [
            {
                "from": "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5",
                "to": "0xe688b84b23f322a994a53dbf8e15fa82cdb71127",
                "hash": "0x7e13dbdfeaae5c54b565fff2b956ac87b04a8282d275e1b4fa8e800020e5228a",
                "value": "0x7880111bfa36f6",
                "type": "outbound"
            }
        ]
    }
    ```

3. GET /current-block

    Example:

    ```
        curl --request GET  --url http://localhost:8787/current-block 
    ```
    Response:
    ```
        {
	        "currentBlock": 20220631
        }
    ```

Assumptions:
1. Recent transactions are more important. Transactions since genesis block are not needed, only grabbing transactions of 10 blocks before starting the service.
2. Transactions can be stored and shown in getTransaction endpoint in ascending order
3. Running parser on interval of 10 seconds is fine.
4. Not using websockets due to constraint of avoiding external libraries and different node url


Potential improvements:
1. Move nodeUrl and parseInterval to environment variables
2. Run server with websockets to listen to latest block instead of running in intervals
3. Create a worked who runs every so often parsing from current block backwards to genesis (if needed) filling historical transactions per address
