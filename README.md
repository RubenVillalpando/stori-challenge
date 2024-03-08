# stori-challenge

Challenge for backend position at stori

# Code rundown

The code for the service is divided in handlers, databases and models and it is initialized in the app.go

The 2 main programs are the report-lambda and the txn-service

The txn-service handles all the "CRUD" operations on managing users, accounts and transactions.

The "lambda" is the one responsible for sending the email off of the csv file.

# How to run the project

clone the project

if you haven't already, [install go](https://go.dev/doc/install)

    git clone https://github.com/RubenVillalpando/stori-challenge.git

use the different make commands to spin up different things
