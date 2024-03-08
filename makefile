ECR_REPOSITORY_URI = 533267442883.dkr.ecr.us-east-2.amazonaws.com/transaction-repository

# GO
run-server:
	cd txn-service && go run cmd/app/main.go
run-lambda:
	cd report-lambda && go run .

# CDK
synth-service:
	cd cdk-service && cdk synth 
deploy-service:
	cd cdk-service && cdk deploy 

synth-lambda:
	cd cdk-lambda && cdk synth
deploy-lambda:
	cd cdk-lambda && cdk deploy

# DOCKER
mysql-up:
	docker-compose up
docker-image:
	cd txn-service && docker build -t txn-service .

docker-run:
	cd txn-service && docker run -p 8080:8080 txn-service

docker-push:
	aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin 533267442883.dkr.ecr.us-east-2.amazonaws.com
	docker build -t $(ECR_REPOSITORY_URI):latest txn-service
	docker push $(ECR_REPOSITORY_URI):latest
