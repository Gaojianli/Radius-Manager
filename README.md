# 🌐 RADIUS Manager

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker Pulls](https://img.shields.io/docker/pulls/gaojianli2333/radius-manager)](https://hub.docker.com/r/gaojianli2333/radius-manager)
[![Docker Image Size](https://img.shields.io/docker/image-size/gaojianli2333/radius-manager/latest)](https://hub.docker.com/r/gaojianli2333/radius-manager)
[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-green.svg)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.3+-blue.svg)](https://www.typescriptlang.org/)

A modern RADIUS user management system built with **CloudWeGo Hertz** and **Vue 3**. Provides complete user authentication, authorization, and management features for FreeRADIUS with a responsive web interface.

The codes are all written by Claude code.

[中文文档](./README.CN.md)

## ✨ Project Highlights

- 🎨 **Modern UI**: Responsive interface based on Arco Design, perfect support for desktop and mobile
- 🔐 **Security & Reliability**: SHA256+Salt password encryption, JWT authentication, role-based access control
- ⚡ **High Performance**: CloudWeGo Hertz high-performance framework with concurrent processing support
- 🌍 **Standard Compatible**: Fully compatible with FreeRADIUS authentication flow and RADIUS protocol
- 📦 **Containerized Deployment**: Docker Compose one-click deployment, ready to use out of the box
- 📱 **Mobile Friendly**: Complete mobile adaptation with touch operations and responsive layout

## 🚀 Features

### 🖥️ **Modern Web Management Interface**
- ✨ Built with Vue 3 + TypeScript + Arco Design for modern interface
- 📱 Complete mobile adaptation with touch operations and responsive layout
- 🎨 Intuitive and user-friendly interface, simple and easy to understand
- ⚡ Real-time data display and operation feedback with hot updates

### 🔧 **Powerful Backend Features**
- 🔐 Secure user authentication and password management system
- 👥 Complete user lifecycle management (create, edit, ban, delete)
- 🛡️ Role-based fine-grained permission control system
- 🔒 SHA256+Salt password encryption following security best practices
- 🎯 JWT stateless authentication supporting distributed deployment
- 🔌 RESTful API design, easy to integrate and extend
- 📊 Detailed authentication logs and statistical analysis
- 🌐 Fully compatible with FreeRADIUS authentication flow
- 📦 Docker Compose one-click deployment supporting containerized environments

## 💻 Tech Stack

### 🎨 Frontend Tech Stack
- **Core Framework**: Vue 3.4+ + TypeScript 5.3+
- **UI Component Library**: Arco Design 2.56+ (officially recommended enterprise-level component library)
- **State Management**: Pinia (next-generation Vuex)
- **Router Management**: Vue Router 4.2+
- **HTTP Requests**: Axios 1.6+ (supports interceptors and automatic retry)
- **Development Tools**: Vite 5.0+ (lightning-fast hot reload and building)
- **Code Quality**: TypeScript type checking + Unplugin auto-import

### ⚙️ Backend Tech Stack
- **High-Performance Framework**: CloudWeGo Hertz 0.10+ (ByteDance's open-source enterprise-level framework)
- **Data Storage**: MySQL 8.0+ (supports transactions and concurrency control)
- **ORM Framework**: GORM 1.30+ (powerful Go ORM)
- **Authentication & Authorization**: JWT (HMAC-SHA256 signature + automatic expiration refresh)
- **Password Security**: SHA256 + 32-byte random salt encryption
- **Deployment**: Docker & Docker Compose + Makefile automation
- **Cross-Origin Support**: CORS middleware + RESTful API design

## 🚀 Quick Start

### ⚡ Quick Start with Docker

Get up and running in 30 seconds:

```bash
# Download and start the complete setup
curl -fsSL https://raw.githubusercontent.com/Gaojianli/radius-manager/main/docker-compose.yml -o docker-compose.yml && docker-compose up -d
```

Then visit http://localhost:8080 and login with `admin` / `admin123`

### 📦 Using Docker Compose (Recommended)

1. Create docker-compose.yml file
```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: radius_mysql
    environment:
      MYSQL_ROOT_PASSWORD: radius123
      MYSQL_DATABASE: radius_mgnt
      MYSQL_USER: radius
      MYSQL_PASSWORD: radius123
    volumes:
      - mysql_data:/var/lib/mysql
    command: --default-authentication-plugin=mysql_native_password
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-pradius123"]
      timeout: 20s
      retries: 10
      interval: 10s
      start_period: 40s
    networks:
      - radius_net

  radius_mgnt:
    image: gaojianli2333/radius-manager:latest
    container_name: radius_mgnt_app
    environment:
      # Database configuration
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: radius
      DB_PASSWORD: radius123
      DB_NAME: radius_mgnt
      
      # Application configuration
      JWT_SECRET: your-super-secret-jwt-key-change-in-production
      SERVER_PORT: :8080
      
      # Default admin user configuration (customize as needed)
      DEFAULT_ADMIN_USER: admin
      DEFAULT_ADMIN_PASSWORD: admin123
      DEFAULT_ADMIN_EMAIL: admin@example.com
    ports:
      - "8080:8080"
    depends_on:
      mysql:
        condition: service_healthy
    networks:
      - radius_net
    restart: unless-stopped

volumes:
  mysql_data:

networks:
  radius_net:
    driver: bridge
```

2. Start services
```bash
docker-compose up -d
```

3. Access web interface
- Management UI: http://localhost:8080
- API Documentation: http://localhost:8080/api/v1

Default admin account (customizable via environment variables):
- Username: `admin` (or `DEFAULT_ADMIN_USER`)
- Password: `admin123` (or `DEFAULT_ADMIN_PASSWORD`)

### 🔒 Production Deployment

For production use, create a `.env` file or modify the environment variables:

```yaml
environment:
  # Security: Change these values in production!
  JWT_SECRET: "your-production-jwt-secret-key-min-32-chars"
  DEFAULT_ADMIN_PASSWORD: "YourSecurePassword123!"
  
  # Optional: Customize admin user
  DEFAULT_ADMIN_USER: "admin"
  DEFAULT_ADMIN_EMAIL: "admin@yourcompany.com"
```

### 🛠️ Development Environment Setup

For development or building from source:

1. Clone the repository
```bash
git clone https://github.com/Gaojianli/radius-manager.git
cd radius-manager
```

2. Install dependencies
```bash
make install
```

3. Configure environment variables
```bash
cp .env.example .env
# Edit .env file to set database connection information
```

4. Start database (using Docker)
```bash
docker-compose up -d mysql
```

5. Build and run
```bash
# Build frontend
make build-frontend

# Run complete application
make run

# Or start frontend and backend separately for development
make dev-frontend  # Start frontend dev server (port 3000)
make dev-backend   # Start backend dev server (port 8080)
```

6. Build your own Docker image (optional)
```bash
# Build custom image
make docker-build

# Or manually build
docker build -t your-registry/radius-manager .
```

## 🎯 Web Interface Features

### 📱 Mobile Optimization
- **Responsive Design**: Adaptive to desktop, tablet, and mobile screens
- **Touch Optimization**: Support for touch gestures, click feedback, and swipe operations
- **Mobile Navigation**: Sidebar overlay design with gesture control and mask layer
- **Mobile Tables**: Smart column hiding, horizontal scrolling, and touch optimization
- **Floating Actions**: FAB (Floating Action Button) design following Material Design guidelines
- **Modal Adaptation**: Mobile modal auto-adapts to screen size

### 🔑 Login Page
- 🔒 Secure user authentication and session management
- ❌ Friendly error messages and form validation
- 💾 Automatic login state persistence with auto-login support
- 📱 Perfect mobile adaptation and touch optimization

### 📊 Dashboard
- 📈 Real-time system statistics and data visualization
- 👥 User activity and authorization statistics (supports display by role)
- ⚙️ System information and version display
- 🎨 Responsive card layout adapting to screen size

### 👥 User Management (Admin)
- 📋 Paginated user list with search and filtering support
- ➕ Visual user creation form with real-time validation
- 🔐 One-click user password reset (admin cannot see passwords)
- 🚫 Flexible user ban/unban functionality
- 🗑️ Safe user deletion (with confirmation dialogs)
- 🏷️ Real-time user status indicators and role management
- 🚀 Mobile FAB floating button supporting manual operations

### 👤 User Profile
- 🗺️ Detailed personal information display (username, email, role, status)
- 🔐 Secure password modification with password strength checking
- 🛡️ Built-in security tips and password security recommendations
- 📅 Creation time and update time tracking

### 📋 Authentication Logs (Admin)
- 📊 Detailed authentication log records and analysis
- 🔍 Support for filtering and searching by username
- 📅 Time range queries and sorting
- 📱 Mobile-adapted table display

## 📡 API Documentation

### Authentication

#### User Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

Response:
```json
{
  "code": 200,
  "expire": "2024-01-02T15:04:05Z",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
Authorization: Bearer <token>
```

### User Management

#### Get User Profile
```http
GET /api/v1/user/profile
Authorization: Bearer <token>
```

#### Change Password
```http
PUT /api/v1/user/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "oldpass123",
  "new_password": "newpass123"
}
```

### Admin Functions

#### Create User
```http
POST /api/v1/admin/users
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123",
  "is_admin": false
}
```

#### Get All Users
```http
GET /api/v1/admin/users?page=1&limit=20
Authorization: Bearer <admin-token>
```

#### Ban/Unban User
```http
PUT /api/v1/admin/users/{user_id}/ban
Authorization: Bearer <admin-token>
```

#### Delete User
```http
DELETE /api/v1/admin/users/{user_id}
Authorization: Bearer <admin-token>
```

### FreeRADIUS Integration

FreeRADIUS integration uses the standard Authorize → Authenticate two-phase verification process:

#### Authorization Interface
Used to check if a user is authorized to authenticate, without password verification.

```http
POST /api/v1/radius/authorize
Content-Type: application/json

{
  "username": "testuser"
}
```

Response:
```json
{
  "result": "accept",
  "username": "testuser", 
  "message": "User authorized to authenticate",
  "attrs": {
    "Auth-Type": "REST",
    "User-ID": "testuser",
    "User-Role": "user"
  }
}
```

#### Authentication Interface
Used to verify user credentials (username and password).

```http
POST /api/v1/radius/auth
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

Response:
```json
{
  "result": "accept",
  "username": "testuser",
  "message": "Authentication successful",
  "attrs": {
    "User-ID": "testuser",
    "User-Role": "user",
    "Session-Type": "authenticated"
  }
}
```

## 🔧 FreeRADIUS Integration

### Two-Phase Verification Process

FreeRADIUS uses a two-phase verification mode to ensure security and performance:

1. **Authorize Phase**
   - Check if user exists
   - Verify account status (disabled, banned)
   - Set authentication type
   - **No password verification**, quickly filter invalid users

2. **Authenticate Phase**
   - Verify user credentials (password)
   - Only execute for users who pass authorization
   - Return final authentication result

### Advantages

- **Performance Optimization**: Invalid users are blocked in authorization phase, reducing password verification overhead
- **Security Enhancement**: Banned users cannot enter authentication phase
- **Standard Compliant**: Follows RADIUS protocol design philosophy
- **Easy to Monitor**: Can separately track authorization and authentication success/failure rates

## 📋 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| **Database Configuration** | | |
| DB_HOST | localhost | MySQL host address |
| DB_PORT | 3306 | MySQL port |
| DB_USER | root | MySQL username |
| DB_PASSWORD | - | MySQL password |
| DB_NAME | radius_mgnt | Database name |
| **Application Configuration** | | |
| JWT_SECRET | your-secret-key | JWT signing key |
| SERVER_PORT | :8080 | Server port |
| **Default Admin Configuration** | | |
| DEFAULT_ADMIN_USER | admin | Default admin username (created only if no admin exists) |
| DEFAULT_ADMIN_PASSWORD | admin123 | Default admin password (created only if no admin exists) |
| DEFAULT_ADMIN_EMAIL | admin@example.com | Default admin email (created only if no admin exists) |

### 🔐 Security Notes

- **Change default credentials**: Always modify `DEFAULT_ADMIN_PASSWORD` and `JWT_SECRET` in production
- **One-time creation**: Default admin is created only when no admin users exist in the database
- **Environment override**: Use environment variables to customize admin credentials without code changes

## 🔒 Security Features

- **Password Encryption**: SHA256 + 32-byte random salt
- **JWT Authentication**: HMAC-SHA256 signature
- **Access Control**: Role-based access control
- **Password Policy**: Minimum 6 characters
- **API Security**: All sensitive operations require authentication
- **Frontend Route Guards**: Unauthenticated users auto-redirect to login
- **Auto Login Expiration**: Auto logout on token expiration

## 🏗️ Project Structure

```
radius-manager/
├── web/                    # Vue3 frontend project
│   ├── src/
│   │   ├── components/     # Reusable components
│   │   ├── views/          # Page components
│   │   ├── layouts/        # Layout components
│   │   ├── stores/         # Pinia state management
│   │   ├── services/       # API services
│   │   └── router/         # Router configuration
│   ├── package.json
│   └── vite.config.ts
├── static/                 # Static file server
│   ├── dist/               # Frontend build output
│   └── static.go           # Go embed static files
├── config/                 # Configuration files
├── controllers/            # Controllers
├── dao/                    # Data access layer
├── database/               # Database connection
├── middleware/             # Middlewares
├── models/                 # Data models
├── routes/                 # Route configuration
├── scripts/                # Build scripts
├── docker-compose.yml      # Docker orchestration
├── Dockerfile              # Container build
├── Makefile               # Build commands
└── README.md              # Project documentation
```

## 🛠️ Development Commands

```bash
make help              # Show all available commands
make install          # Install frontend and backend dependencies
make build            # Build entire project
make build-frontend   # Build frontend only
make build-backend    # Build backend only
make run              # Build and run project
make dev-frontend     # Start frontend dev server
make dev-backend      # Start backend dev server
make clean            # Clean build files
make docker-build     # Build Docker image
make docker-run       # Start with Docker Compose
make test             # Run tests
make format           # Format code
make lint             # Code linting
```

## 🐛 Troubleshooting

### Frontend Won't Load

1. Ensure frontend is built: `make build-frontend`
2. Check if `static/dist/` directory has files
3. Confirm server port configuration is correct

### API Request Failures

1. Check if backend service is running on correct port
2. Verify JWT token is valid
3. Check browser network panel for error messages

### Database Connection Issues

1. Confirm MySQL service status
2. Check database configuration in `.env` file
3. Verify database user permissions

### FreeRADIUS Authentication Failures

1. Check REST module configuration is correct
2. Verify API endpoint accessibility
3. View FreeRADIUS debug output: `sudo freeradius -X`

## 🤝 Contributing

We welcome contributions of all kinds! Whether it's bug reports, feature requests, documentation improvements, or code contributions.

### 📋 Contribution Process

1. **Fork the project**: Click the Fork button in the top right
2. **Clone locally**: `git clone https://github.com/YOUR_USERNAME/radius-manager.git`
3. **Create feature branch**: `git checkout -b feature/amazing-feature`
4. **Local development**:
   ```bash
   # Install dependencies
   make install
   
   # Start development environment
   make dev-frontend    # Frontend dev server
   make dev-backend     # Backend dev server
   ```
5. **Code standards**:
   ```bash
   make format         # Format code
   make lint          # Code linting
   make test          # Run tests
   ```
6. **Commit changes**:
   ```bash
   git add .
   git commit -m "✨ Add amazing feature: feature description"
   git push origin feature/amazing-feature
   ```
7. **Create Pull Request**: Go back to GitHub to create PR

### 🎯 Contribution Types

- 🐛 **Bug Fixes**: Fix issues with existing functionality
- ✨ **New Features**: Add new features
- 📚 **Documentation**: Improve documentation and instructions
- 🎨 **UI/UX**: Interface and user experience improvements
- ⚡ **Performance**: Performance optimization and improvements
- 🔧 **Refactoring**: Code refactoring and architecture optimization

### 📝 Code Standards

- **Go Code**: Follow `gofmt` and `golint` standards
- **Vue Code**: Use TypeScript and Vue 3 Composition API style
- **Commit Messages**: Use [Conventional Commits](https://www.conventionalcommits.org/) format
- **Documentation**: All public APIs must have corresponding documentation

## 🌟 Acknowledgments

Thanks to all developers who contributed to this project!

## 📄 License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).

```
MIT License

Copyright (c) 2024 RADIUS Manager Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## 📞 Contact

- 📧 **Issues**: [GitHub Issues](https://github.com/Gaojianli/radius-manager/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/Gaojianli/radius-manager/discussions)
- 🐛 **Bug Reports**: [Bug Report Template](https://github.com/Gaojianli/radius-manager/issues/new?template=bug_report.md)
- ✨ **Feature Requests**: [Feature Request Template](https://github.com/Gaojianli/radius-manager/issues/new?template=feature_request.md)

---

⭐ If this project helps you, please give us a **Star**!

**Happy Coding! 🎉**