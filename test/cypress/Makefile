run-server:
	cd ../../ && \
	sed -i "s/TEST_CLIENT_BACKEND_URL=localhost:8000/TEST_CLIENT_BACKEND_URL=localhost:6002/g" server/.env.dev && \
	sed -i "s/TEST_CLIENT_BACKEND_URI=http:\/\/localhost:8000/TEST_CLIENT_BACKEND_URI=http:\/\/localhost:6002/g" server/.env.dev && \
	make setup_and_run

run-wgp-backend:
	cd ../../sp_mocks/wgp/backend && \
	echo "PORT=6002" > .env && \
	echo "FRONTEND_URL=http://localhost:5173" >> .env && \
	echo "LAYER8_URL=http://localhost:5001" >> .env && \
	echo "BACKEND_URL=http://localhost:6002" >> .env && \
	npm install && node server.js

run-wgp-frontend:
	cd ../../sp_mocks/wgp/frontend && \
	npm install && npm run dev
