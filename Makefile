## NPM Install sp_mocks
npm_install_wgp:
	cd sp_mocks/wgp/frontend && npm install && npm i layer8_interceptor && cd ../backend && npm install && npm i layer8_middleware

npm_install_imsharer:
	cd sp_mocks/imsharer/frontend && npm install && npm i layer8_interceptor && cd ../backend && npm install && npm i layer8_middleware

npm_install_all:
	make npm_install_wgp && make npm_install_imsharer

## Build bundled .js file with bip-39 & bip-32 related functions
build_mnemonic_bundle:
	cd server/mnemonic && npm install && npm run build

build_scram_bundle:
	cd server/scram && npm install && npm run build

test_scram_ts:
	cd server/scram && npm run test

## Run Service Provider Mocks
run_wgp_frontend: # Port 5173
	cd sp_mocks/wgp/frontend && npm run dev
	
run_wgp_backend: # Port 8000
	cd sp_mocks/wgp/backend && npm run dev

run_imsharer_frontend:
	cd sp_mocks/imsharer/frontend && npm run dev
	
run_imsharer_backend:
	cd sp_mocks/imsharer/backend && npm run dev


# Build and Run Resource Server and Proxy
go_mod_tidy:
	cd ./server && go mod tidy

go_test:
	cd server && go test ./... -v -cover

run_server: # Port 5001
	cd server/cmd/app && go run main.go

# Build and Push Docker Images
build_server_image:
	docker build --tag layer8-server-new --file Dockerfile .

build_sp_mocks_frontend_image:
	cd sp_mocks/wgp/frontend && docker build --tag sp_mocks_frontend --file Dockerfile .

build_sp_mocks_backend_image:
	cd sp_mocks/wgp/backend && docker build --tag sp_mocks_backend --file Dockerfile .


# To build all images at once
build_images:
	make build_server_image && make build_sp_mocks_frontend_image && make build_sp_mocks_backend_image

run_layer8_server_image:
	docker run -p 5001:5001 -t layer8-server-new

run_sp_mocks_frontend_image:
	docker run -p 8080:8080 -t sp_mocks_frontend

run_sp_mocks_backend_image:
	docker run -p 8000:8000 -t sp_mocks_backend

push_layer8_server_image:
	aws lightsail push-container-image --region ca-central-1 --service-name aws-container-service-t1 --label layer8-server-version-stable-9 --image layer8-server-new:latest

push_sp_mocks_frontend_image:
	aws lightsail push-container-image --region ca-central-1 --service-name container-service-2 --label frontendversiontest2 --image sp_mocks_frontend:latest

push_sp_mocks_backend_image:
	aws lightsail push-container-image --region ca-central-1 --service-name container-service-3 --label backendtest2 --image sp_mocks_backend:latest

push_images:
	make push_layer8_server_image && make push_sp_mocks_frontend_image && make push_sp_mocks_backend_image


# Run a local Postgres DB
run_local_db:
	docker run -d --rm \
		--name layer8-resource \
		-v $(PWD)/.docker/postgres:/var/lib/postgresql/data \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DBNAME=postgres \
		-p 5434:5432 postgres:14.3


setup_local_dependency:
	cd server && go run cmd/setup/setup.go

run_layer8server_local:
	cd server && go run cmd/app/main.go

setup_and_run: 
	make setup_local_dependency && make run_layer8server_local

mockgen:
	mockgen -source=server/internals/service/service.go -destination=server/utils/mocks/internal_service_mock.go -package=mocks

SP_MOCK := wgp
set_client_creds:
	db_name=$$(cat server/.env | grep ^DB_NAME | cut -d '=' -f2); \
	db_user=$$(cat server/.env | grep ^DB_USER | cut -d '=' -f2); \
	client_username=$$(cat server/.env | grep TEST_CLIENT_USERNAME | cut -d '=' -f2); \
	client_id=$$(docker exec layer8-postgres psql -U $$db_user -d $$db_name \
		-c "SELECT id FROM clients WHERE username='$$client_username'" -t -A); \
	client_secret=$$(docker exec layer8-postgres psql -U $$db_user -d $$db_name \
		-c "SELECT secret FROM clients WHERE username='$$client_username'" -t -A); \
	if [ -z "$$client_id" ]; then \
		echo "Client not found"; \
	else \
		if [ "$(SP_MOCK)" = "wgp" ]; then \
			echo "LAYER8_CLIENT_ID=$$client_id" >> sp_mocks/wgp/backend/.env; \
			echo "LAYER8_CLIENT_SECRET=$$client_secret" >> sp_mocks/wgp/backend/.env; \
		elif [ "$(SP_MOCK)" = "imsharer" ]; then \
			echo "LAYER8_CLIENT_ID=$$client_id" >> sp_mocks/imsharer/backend/.env; \
			echo "LAYER8_CLIENT_SECRET=$$client_secret" >> sp_mocks/imsharer/backend/.env; \
		else \
			echo "Invalid SP_MOCK"; \
		fi; \
	fi
