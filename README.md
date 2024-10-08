# Blog App <span style="color: #00ADD8;">Go</span>

## Overview

This is my first project in Go, where I've developed a fully functional blog application with a focus on security, efficiency, and user experience. Key features and technologies implemented include:

- `JWT Authentication`: Secure user authentication using JSON Web Tokens for session management.

- `Middlewares`: Developed custom middleware to verify JWTs, ensuring only authenticated users can access protected routes.

- `Password Hashing`: Leveraged bcrypt for secure password storage, ensuring user credentials are safely hashed.

- `OTP Verification`: Implemented a robust signup process with One-Time Password (OTP) verification, enhancing user registration security.

- `Cloudinary Integration`: Successfully integrated Cloudinary for image uploads, allowing users to manage their profile pictures seamlessly.

- `Database Interaction without ORM`: Utilized raw SQL queries for MySQL database interactions, gaining hands-on experience with direct database management and optimization.

- `Input Validation`: Incorporated comprehensive input validation to ensure data integrity and prevent malicious input.

- `Bloom Filters`: Implemented Bloom filters to efficiently check for existing emails during user registration, optimizing performance for large databases.

- `CORS Management`: Configured Cross-Origin Resource Sharing (CORS) to enable secure interactions between the frontend and backend.

This project not only showcases my technical skills in Go but also reflects my commitment to building secure and efficient applications. The experience gained from this project has strengthened my understanding of backend development principles and has equipped me with practical knowledge in handling real-world application challenges.

## Tech Used

- `Go` Built the HTTP server and backend logic.
- `Cloudinary` Cloud service to handle image uploads for user profiles.
- `Docker` Containerized the application for ease of deployment and consistency across environments.

#### <span style="color: #68217A">As part of my learning journey, I'm excited to experiment with building this blog application in C# and PostgreSQL to explore new technologies and enhance my skills!</span>

## Set-Up the Project locally

### Go

1. Clone the repository to the local:

   ```bash
   git clone https://github.com/Sahil2k07/Blog-App-Go
   ```

2. Move to the project directory:

   ```bash
   cd Blog-App-Go
   ```

3. Set up all the required env variable by making a `.env` file. A `.env.example` file has been given for reference.

   ```dotenv
   PORT=localhost:3000

   JWT_SECRET=GoAppSecret

   ALLOWED_ORIGINS=*

   MYSQL_URL=root:YOUR_PASSWORD@tcp(127.0.0.1:3306)/Blog_App_Go

   # Mailer Details.
   MAIL_HOST=
   MAIL_USER=
   MAIL_PASS=

   # Cloudinary Details.
   CLOUD_NAME=
   API_KEY=
   API_SECRET=
   FOLDER_NAME=Blog_App_Go
   ```

4. Run the command to download all the dependencies to your local machine:

   ```bash
   go mod vendor
   ```

5. Now you will need to apply all the migrations in your database. First install goose locally:

   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

6. Move to the migrations directory:

   ```bash
   cd src/database/migrations
   ```

7. Run the command to apply the migrations. Make sure to modify the command to have your's used `MySQL` database url in the `.env`. Make sure to have a database created before-hand:

   ```bash
   goose mysql "root:YOUR_PASSWORD@tcp(127.0.0.1:3306)/Blog_App_Go" up
   ```

8. After applying the migrations traverse back to the root directory:

   ```bash
   cd ../../../
   ```

9. Build the Binary to start the Project:

   ```bash
   go build
   ```

10. Start the Project:

    ```bash
    ./Blog-App-Go
    ```

### Docker

1. First clone the Project locally:

   ```bash
   git clone https://github.com/Sahil2k07/Blog-App-Go
   ```

2. Move to the Project directory:

   ```bash
   cd Blog-App-Go
   ```

3. Set these environment variables in the `.env` file.

   ```dotenv
   # Mailer Details.
   MAIL_HOST=
   MAIL_USER=
   MAIL_PASS=

   # Cloudinary Details.
   CLOUD_NAME=
   API_KEY=
   API_SECRET=
   FOLDER_NAME=Blog_App_Go
   ```

4. Run the command to start your Containerized Application

   ```bash
   docker-compose up
   ```

   or

   ```bash
   docker-compose up -d
   ```

5. If you have Docker Compose Plugin, Use this command instead

   ```bash
   docker compose up
   ```

   or

   ```bash
   docker compose up -d
   ```

6. You will be able to access this application in `localhost:5000` of your machine.
