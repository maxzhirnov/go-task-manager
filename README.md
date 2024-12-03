## **Task Manager**

A secure and user-friendly **Task Manager** application built with **Go** for the backend and a simple HTML/CSS/JavaScript frontend. This project demonstrates authentication with **JWT tokens** (including refresh tokens), user-specific tasks, and includes comprehensive test coverage.

### **Live Demo**
You can try out the application at: [https://todo.mzhirnov.com](https://todo.mzhirnov.com)

#### Test Credentials
- Username: `testuser@temporary.com`
- Password: `123456`

---

### **Features**
1. **User Authentication**:
   - Secure registration and login using hashed passwords (bcrypt).
   - JWT-based authentication with access and refresh tokens.

2. **User-Specific Task Management**:
   - Users can create, view, update, and delete their own tasks.
   - Tasks are bound to user accounts, ensuring data privacy.

3. **Token Refresh Mechanism**:
   - Automatic token renewal using refresh tokens for seamless user experience.

4. **Frontend**:
   - Simple and clean UI for managing tasks.
   - Dynamic task loading, creation, and deletion.

5. **Backend**:
   - RESTful API built with **Go**.
   - Database integration using PostgreSQL.

6. **Dockerized Deployment**:
   - Fully containerized for easy deployment with Docker and Docker Compose.

7. **Comprehensive Testing**:
   - Unit tests for models, handlers, and middleware ensure reliability.

---

### **Technologies Used**
- **Backend**:
  - Go (Golang)
  - Gorilla Mux (Router)
  - bcrypt (Password hashing)
  - jwt-go (JSON Web Tokens)
- **Frontend**:
  - HTML/CSS/JavaScript
- **Database**:
  - PostgreSQL
- **Containerization**:
  - Docker & Docker Compose
- **Testing**:
  - sqlmock (Mock database for testing)
  - testify (Assertions in tests)

---

### **Getting Started**

#### **Prerequisites**
- Docker and Docker Compose installed on your machine.
- Go (Golang) installed (optional for development).

#### **Clone the Repository**
```bash
git clone https://github.com/maxzhirnov/go-task-manager.git
cd go-task-manager
```

---

### **Setup and Run the Application**

#### **1. Run with Docker Compose**
The easiest way to run the application is using Docker Compose. It will set up the backend, database, and serve the frontend.

```bash
docker-compose up --build
```

- The backend will be available at: `http://localhost:8080`
- The frontend can be accessed via the same address.

---

#### **2. Run Locally**

If you want to run the application locally without Docker:

1. **Setup PostgreSQL**:
   - Ensure PostgreSQL is running locally.
   - Create a database named `taskmanager`.
   - Run the SQL scripts in `scripts/db/init.sql` to set up the schema.

2. **Set Environment Variables**:
   Create a `.env` file or export the variables directly:
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=mysecretpassword
   export DB_NAME=taskmanager
   export PORT=8080
   ```

3. **Run the Backend**:
   ```bash
   go run cmd/api/main.go
   ```

4. **Serve the Frontend**:
   Open the `web/index.html` file in your browser.

---

### **Endpoints**

#### **Authentication**
| Method | Endpoint         | Description                |
|--------|------------------|----------------------------|
| POST   | `/api/register`  | Register a new user        |
| POST   | `/api/login`     | Login and get tokens       |
| POST   | `/api/refresh`   | Get a new access token     |

#### **Tasks**
| Method | Endpoint         | Description                |
|--------|------------------|----------------------------|
| GET    | `/api/tasks`     | Get all tasks for a user   |
| POST   | `/api/tasks`     | Create a new task          |
| GET    | `/api/tasks/{id}`| Get details of a task      |
| PUT    | `/api/tasks/{id}`| Update a specific task     |
| DELETE | `/api/tasks/{id}`| Delete a specific task     |

---

### **Sample `.env` File**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=mysecretpassword
DB_NAME=taskmanager
PORT=8080
```

---

### **Testing**

#### **Run Unit Tests**
```bash
go test ./... -v
```

#### **Run with Coverage**
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

### **Screenshots**

#### **Login Page**
Login Page

#### **Task Management Page**
Task Management Page

---

### **Roadmap**
1. Add pagination for tasks.
2. Implement role-based access control (e.g., admin vs. user).
3. Add more comprehensive integration tests.
4. Deploy the application to a cloud platform.

---

### **License**

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

### **Contributing**

Contributions are welcome! Please open an issue or submit a pull request to contribute to this project.