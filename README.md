# MiniApp

MiniApp is a web application with a **Golang backend**, a **Next.js frontend**, and a **PostgreSQL database**. This document explains how to set up, build, test, and deploy the application.

---

## **Prerequisites**

Before starting, ensure you have the following installed:
- Docker and Docker Compose
- PostgreSQL
- Node.js (v18 or later) and NPM
- Golang (v1.20 or later)
- NGINX (for deployment)
- Certbot (for SSL, optional)

---

## **1. Installation with Docker**

### **Step 1: Clone the Repository**
```bash
git clone https://github.com/your-repo/miniapp.git
cd miniapp
```

### **Step 2: Build and Run with Docker Compose**
Ensure `docker-compose.yml` is present in the root directory.

Run the following command to start all services:
```bash
docker-compose up --build
```

Services:
- **Backend**: `http://localhost:8080`
- **Frontend**: `http://localhost:3000`
- **PostgreSQL**: `localhost:5432`

---

## **2. Manual Setup (Without Docker)**

### **Step 1: Set Up PostgreSQL**
Run the following SQL commands to create the database and user:
```sql
CREATE USER miniapp_user WITH PASSWORD 'miniapp_user_12345';
CREATE DATABASE mini_db;
GRANT ALL PRIVILEGES ON DATABASE mini_db TO miniapp_user;
```

Update the `.env` file in the `backend` directory with your database credentials:
```plaintext
DB_HOST=localhost
DB_NAME=mini_db
DB_USER=miniapp_user
DB_PASSWORD=miniapp_user_12345
DB_SSLMODE=disable
DB_PORT=5432
SECRET=9y/scCRS5dLC7HvpPHQKYk+OINUdPSYKcjABnmgcxns=
```

---

### **Step 2: Build and Run the Backend**
Navigate to the `backend` directory and run:
```bash
cd backend
go mod tidy
go build -o main
./main
```

The backend will start on `http://localhost:8080`.

---

### **Step 3: Build and Run the Frontend**
Navigate to the `frontend` directory and run:
```bash
cd frontend
npm install
npm run build
npm start
```

The frontend will start on `http://localhost:3000`.

---

## **3. Deploy with NGINX and Certbot**

### **Step 1: Configure NGINX**
Copy the provided `nginx` configuration file to the NGINX sites-available directory:
```bash
sudo cp nginx.conf /etc/nginx/sites-available/default
sudo nginx -t
sudo systemctl restart nginx
```

**NGINX Configuration:**
```nginx
server {
    server_name miniapp.dandanjan.ir;

    location / {
        proxy_pass http://127.0.0.1:3000;
    }

    location /backend {
        proxy_pass http://127.0.0.1:8080;
    }

    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/miniapp.dandanjan.ir/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/miniapp.dandanjan.ir/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot
}

server {
    if ($host = miniapp.dandanjan.ir) {
        return 301 https://$host$request_uri;
    } # managed by Certbot

    listen 80;
    server_name miniapp.dandanjan.ir;
    return 404; # managed by Certbot
}
```

### **Step 2: Set Up SSL with Certbot**
Install Certbot and generate SSL certificates:
```bash
sudo apt update
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d miniapp.dandanjan.ir
```

Test the SSL configuration:
```bash
sudo nginx -t
sudo systemctl reload nginx
```

---

## **4. Running Tests**

### **Backend Tests**
Navigate to the `backend` directory and run:
```bash
cd backend
go test ./... -v
```

### **Frontend Tests**
Navigate to the `frontend` directory and run:
```bash
cd frontend
npm install
npm run test
```

---

## **5. API Documentation with Swagger**

### **Generate and Serve Swagger API Documentation**
1. Place your `swagger.yml` in the `backend` directory.
2. Use a tool like `swagger-ui` to serve the documentation:
   ```bash
   docker run -p 8081:8080 -v $(pwd)/swagger.yml:/usr/share/nginx/html/swagger.yml swaggerapi/swagger-ui
   ```
3. Open the Swagger UI in your browser at `http://localhost:8081`.

---

## **6. Environment Variables**

### Backend
```plaintext
DB_HOST=localhost
DB_NAME=mini_db
DB_USER=miniapp_user
DB_PASSWORD=miniapp_user_12345
DB_SSLMODE=disable
DB_PORT=5432
SECRET=9y/scCRS5dLC7HvpPHQKYk+OINUdPSYKcjABnmgcxns=
```

### Frontend
```plaintext
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## **7. Useful Commands**

### Restart Docker Compose
```bash
docker-compose down
docker-compose up --build
```

### Restart NGINX
```bash
sudo systemctl restart nginx
```

### Logs
- Backend logs are available in the terminal where the backend runs or in Docker logs.
- Frontend logs are available in the browser console or Docker logs.

---

## **8. Directory Structure**
```plaintext
miniapp/
├── backend/            # Golang backend code
├── frontend/           # Next.js frontend code
├── docker-compose.yml  # Docker Compose configuration
├── fullchain.pem       # SSL certificate (for HTTPS)
├── nginx.conf          # NGINX configuration file
├── privkey.pem         # SSL private key (for HTTPS)
├── README.md           # Documentation file
└── swagger.yaml        # API documentation file
```
