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

use the different make commands to spin up different things, mainly:

    make mysql-up

to spin up a docker mysql instance and to run each go program

    make run-server
    make run-lambda

# Notes to add

I attempted to deploy via cdk which was successful but had issues connecting with the mysql instance in the cloud as well as the s3 bucket so everything is now done locally :((((

Thanks for taking the time to look at my project
