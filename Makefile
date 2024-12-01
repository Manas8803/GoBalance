.PHONY: build deploy clean all test

build:
	GOOS=linux GOARCH=amd64 go build -o ./deploy_scripts/assets/app_server/app_server ./app_server/main.go  
	GOOS=linux GOARCH=amd64 go build -o ./deploy_scripts/assets/load_balancer/load_balancer ./load_balancer/main.go  
	cd ./deploy_scripts/assets && zip -r app_server.zip app_server/
	cd ./deploy_scripts/assets && zip -r load_balancer.zip load_balancer/

deploy:
	cd deploy_scripts && cdk deploy --require-approval never

destroy:
	cd deploy_scripts && cdk destroy

clean:
	rm -rf ./deploy_scripts/assets/app_server
	rm -rf ./deploy_scripts/assets/load_balancer 
	rm -rf ./deploy_scripts/assets/app_server.zip 
	rm -rf ./deploy_scripts/assets/load_balancer.zip

test:
	./test.sh

all:
	- make destroy
	- make clean
	make build
	make deploy
