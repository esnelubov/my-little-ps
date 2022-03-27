# my-little-ps

A make-believe payment system that allows to create wallets for its clients,
transfer "money" to these wallets or between them and prepare reports about
these operations.

## Requirements
This code was tested with the configuration:
* Ubuntu 20.04
* Go 1.17
* PostgreSQL 12

## Usage

Build the project using:

    make -f ./Makefile.mk build

This will create a `bin` folder with all executables, a `settings.yml` file and migrations folder inside.

By default migrations are applied automatically.

You can apply migrations manually by turning the `autoMigrate` setting off in `settings.yml` file and using https://github.com/golang-migrate/migrate#cli-usage

my-little-ps consists of 3 services:
1. `gateway` - listens for API requests
2. `op_processor` - processes payment operations
3. `wlt_balancer` - an optional service that allows for several `op_processor` instances to effectively process operations

First you should start the `gateway` service

Then you should update currency rates. Send the following message to the gateway:

**POST** to `http://127.0.0.1:4567/update_currencies` with body

    {
        "Rates": {
            "USD": 1.000000,
            "EUR": 1.093809,
            "RUB": 0.009803,
            "THB": 0.029758
        }
    }

Changes to currency rates are updated internally every 5 minutes. To see the changes immediately please restart the gateway.

Normally all messages return `200 OK` response with JSON body that contains the `payload` field. 
If something goes wrong then the `payload` will contain an `error` field with an error description. 

Now you can create some wallets. Send messages like this to the gateway:

**POST** to `http://127.0.0.1:4567/wallet` with body

    {
        "Name":     "Jon Snow",
        "Country":  "UK",
        "City":     "London",
        "Currency": "EUR"
    }

The reply will contain the wallet ID. 

If the reply contains the error `currency is not allowed` it means that currency rates wasn't yet updated - they are updated internally every 5 minutes to reduce the load on a DB.

Knowing the wallet ID you can send some money to the wallet:

**POST** to `http://127.0.0.1:4567/receive_amount` with body

    {
        "WalletId": 1,
        "Amount": 200000000,
        "Currency": "USD"
    }

This message will send 200 USD to the wallet 1. my-little-ps stores balances and currency rates as integer numbers. To convert them to integers it multiplies them by 1,000,000.
This means that to send 2.15 US Dollars to the system you should write them as 2150000. You can change this constant in `common/constants/constants.go`

my-little-ps automatically converts amount to the wallet's currency.

Note that receive_amount will not be processed because we didn't start the `op_processor` service.
You can start a simple instance of the `op_processor` service just by executing it without arguments.

Now operations will be processed.

You can transfer money between wallets with the message like this:

**POST** to `http://127.0.0.1:4567/transfer_amount` with body

    {
        "OriginWalletId": 1,
        "TargetWalletId": 3,
        "Amount": 1000000,
        "Currency": "EUR"
    }

This message will transfer 1 EUR from the wallet 1 to the wallet 3.

To see an operations report for a wallet you can send:

**GET** to `http://127.0.0.1:4567/operations/1?from=2022-02-27T21:08:56Z&to=2022-03-27T21:08:56Z&offset=0&limit=1000`

This request will return operations for the wallet 1.

All parameters are optional. Their defaults are: 

* `from` = time.Now() - 1 month 
* `to` = time.Now()
* `offset` = 0
* `limit` = 1000

You can also download the report as CSV:

**GET** to `http://127.0.0.1:4567/operations/file/1?from=2022-02-27T21:08:56Z&to=2022-03-27T21:08:56Z`

Or just get the total In/Out amounts for the wallet:

**GET** to `http://127.0.0.1:4567/operations/total/1?from=2022-02-27T21:08:56Z&to=2022-03-27T21:08:56Z`

## Advanced usage

By default all operations are processed by a single `op_processor` instance. 
If you run several instances then they will spend their time waiting for lock on a wallet which is not effective.

To solve this problem you can run an `op_processor` instance only for the given wallet group:

    op_processor --number 1

This `op_processor` instance will process operations only for the wallets marked with group 1.

By default, all wallets belong to a group 1. You can set groups to them manually through the DB, but there is a more advanced way to do it.

The `wlt_balancer` service automatically assigns groups to wallets based on how many operations there were for each wallet in the given period of time.

You can start `wlt_balancer` to periodically rebalance wallets between 2 `op_processor` instances like this:

    wlt_balancer --number 2

This way `wlt_balancer` will rebalance wallets between processors every 30 minutes (setting `walletBalancerDelay`) 
based on operations count for each wallet.

For example if we have only 4 wallets and during the last 30 minutes there were these many operations for each of them:

* Wallet 1 - 1234 operations
* Wallet 2 - 5678 operations
* Wallet 3 - 0 operations
* Wallet 4 - 3456 operations

Then after rebalancing wallets 1 and 4 will be assigned to the op_processor 1 (1234 + 3456 = 4690) and wallets 2 and 3 will be assigned to op_processor 2 (5678 + 0 = 5678).
This way both processors will have wallets with roughly the same activity. This reduces the possibility of a situation when one `op_processor` has all the active wallets while the other is idling with wallets that are rarely used.

To partition wallets based on operations count `wlt_balancer` uses the Karmarkar–Karp multiway number partitioning algorithm: https://en.wikipedia.org/wiki/Largest_differencing_method

With the current implementation of the algorithm it can partition 10000000 wallets in 100 groups in 1 minute.
