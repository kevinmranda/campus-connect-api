# Campus Connect API Documentation

## Overview
The **Campus Connect API** is a RESTful API designed to support a campus community platform. It allows users to register, log in, create and view posts with images, browse job listings, manage events, and access class timetables. The API is built with Go (Gin framework) and uses a PostgreSQL database with GORM for data management.

This documentation provides details for all available endpoints, including request formats, response structures, and example payloads, to help frontend developers integrate with the API.

## Base URL
- **Development**: `http://localhost:3000`
- **Production**: (To be provided when deployed)

## Authentication
- Currently, all endpoints are public and do not require authentication.

## Static Files
- Images (e.g., post images, user profile pictures) are stored in the `./Images` folder and served at `/Images`.
- Example: A profile image at `./Images/profile-picture-UID1.jpg` can be accessed via `http://localhost:3000/Images/profile-picture-UID1.jpg`.

## Error Responses
Common error responses include:
- **400 Bad Request**: Invalid input or parameters.
  ```json
  {"error": "Invalid input: ..."}
  ```
- **404 Not Found**: Resource not found.
  ```json
  {"error": "User not found"}
  ```
- **500 Internal Server Error**: Server-side error.
  ```json
  {"error": "Failed to fetch users"}
  ```

## Endpoints

### User Endpoints
Manage user registration, login, and profiles.

#### POST /api/user/register
Register a new user. Profile images are saved as `profile-picture-UID<user_id>.<ext>` (e.g., `profile-picture-UID1.jpg`).

- **Request Body**:
  ```json
  {
    "name": "string (required)",
    "profileImage": "string (base64-encoded image, optional, e.g., data:image/jpeg;base64,/9j/...)",
    "role": "string (optional)",
    "course": "string (optional)",
    "year": "string (optional)",
    "email": "string (unique, required)",
    "password": "string (required)"
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "Student",
    "profileImage": "./Images/profile-picture-UID1.jpg"
  }
  ```
- **Example**:
  ```bash
  curl -X POST http://localhost:3000/api/user/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"secret","role":"Student","course":"Computer Science","year":"3rd","profileImage":"data:image/jpeg;base64,/9j/..."}'
  ```

#### POST /api/user/login
Log in a user.

- **Request Body**:
  ```json
  {
    "email": "string (required)",
    "password": "string (required)"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "message": "Login successful",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "Student",
      "course": "Computer Science",
      "year": "3rd",
      "profileImage": "./Images/profile-picture-UID1.jpg"
    }
  }
  ```
- **Example**:
  ```bash
  curl -X POST http://localhost:3000/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secret"}'
  ```

#### GET /api/users
List all users.

- **Response (200 OK)**:
  ```json
  [
    {
      "id": 1,
      "createdAt": "2025-04-24T10:00:00Z",
      "updatedAt": "2025-04-24T10:00:00Z",
      "deletedAt": null,
      "name": "John Doe",
      "profileImage": "./Images/profile-picture-UID1.jpg",
      "role": "Student",
      "course": "Computer Science",
      "year": "3rd",
      "email": "john@example.com",
      "posts": []
    },
    ...
  ]
  ```
- **Example**:
  ```bash
  curl http://localhost:3000/api/users
  ```

#### GET /api/users/:id
Get a user by ID.

- **Path Parameters**:
  - `id`: User ID (integer)
- **Response (200 OK)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T10:00:00Z",
    "deletedAt": null,
    "name": "John Doe",
    "profileImage": "./Images/profile-picture-UID1.jpg",
    "role": "Student",
    "course": "Computer Science",
    "year": "3rd",
    "email": "john@example.com",
    "posts": []
  }
  ```
- **Example**:
  ```bash
  curl http://localhost:3000/api/users/1
  ```

#### PUT /api/users/:id/update
Update a user. Profile images are saved as `profile-picture-UID<user_id>.<ext>`.

- **Path Parameters**:
  - `id`: User ID (integer)
- **Request Body**:
  ```json
  {
    "name": "string (optional)",
    "profileImage": "string (base64-encoded image, optional)",
    "role": "string (optional)",
    "course": "string (optional)",
    "year": "string (optional)",
    "email": "string (optional)",
    "password": "string (optional)"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T12:00:00Z",
    "deletedAt": null,
    "name": "John Smith",
    "profileImage": "./Images/profile-picture-UID1.jpg",
    "role": "Student",
    "course": "Computer Science",
    "year": "4th",
    "email": "john.smith@example.com",
    "posts": []
  }
  ```
- **Example**:
  ```bash
  curl -X PUT http://localhost:3000/api/users/1/update \
  -H "Content-Type: application/json" \
  -d '{"name":"John Smith","year":"4th","profileImage":"data:image/png;base64,iVBORw0KGgo..."}'
  ```

#### DELETE /api/users/:id/delete
Delete a user and their profile image.

- **Path Parameters**:
  - `id`: User ID (integer)
- **Response (200 OK)**:
  ```json
  {"message": "User deleted successfully"}
  ```
- **Example**:
  ```bash
  curl -X DELETE http://localhost:3000/api/users/1/delete
  ```

### Post Endpoints
Manage campus posts (e.g., announcements, discussions) with associated user details.

#### GET /api/posts
List all posts with user details.

- **Response (200 OK)**:
  ```json
  [
    {
      "id": 1,
      "createdAt": "2025-04-24T10:00:00Z",
      "updatedAt": "2025-04-24T10:00:00Z",
      "deletedAt": null,
      "image": "./Images/post-my-post.jpg",
      "title": "My Post",
      "description": "This is a post",
      "userID": 1,
      "user": {
        "id": 1,
        "name": "John Doe",
        "profileImage": "./Images/profile-picture-UID1.jpg",
        "role": "Student",
        "course": "Computer Science",
        "year": "3rd"
      }
    },
    ...
  ]
  ```
- **Example**:
  ```bash
  curl http://localhost:3000/api/posts
  ```

#### POST /api/posts
Create a new post.

- **Request Body**:
  ```json
  {
    "title": "string (required)",
    "description": "string (optional)",
    "image": "string (base64-encoded image, optional)",
    "userID": 1 (required)
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T10:00:00Z",
    "deletedAt": null,
    "image": "./Images/post-my-post.jpg",
    "title": "My Post",
    "description": "This is a post",
    "userID": 1,
    "user": {
      "id": 1,
      "name": "John Doe",
      "profileImage": "./Images/profile-picture-UID1.jpg",
      "role": "Student",
      "course": "Computer Science",
      "year": "3rd"
    }
  }
  ```
- **Example**:
  ```bash
  curl -X POST http://localhost:3000/api/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"My Post","description":"This is a post","image":"data:image/jpeg;base64,/9j/...","userID":1}'
  ```

#### GET /api/posts/:id
Get a post by ID with user details.

- **Path Parameters**:
  - `id`: Post ID (integer)
- **Response (200 OK)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T10:00:00Z",
    "deletedAt": null,
    "image": "./Images/post-my-post.jpg",
    "title": "My Post",
    "description": "This is a post",
    "userID": 1,
    "user": {
      "id": 1,
      "name": "John Doe",
      "profileImage": "./Images/profile-picture-UID1.jpg",
      "role": "Student",
      "course": "Computer Science",
      "year": "3rd"
    }
  }
  ```
- **Example**:
  ```bash
  curl http://localhost:3000/api/posts/1
  ```

#### PUT /api/posts/:id/update
Update a post.

- **Path Parameters**:
  - `id`: Post ID (integer)
- **Request Body**:
  ```json
  {
    "title": "string (optional)",
    "description": "string (optional)",
    "image": "string (base64-encoded image, optional)",
    "userID": 1 (optional)
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T12:00:00Z",
    "deletedAt": null,
    "image": "./Images/post-updated-post.jpg",
    "title": "Updated Post",
    "description": "Updated description",
    "userID": 1,
    "user": {
      "id": 1,
      "name": "John Doe",
      "profileImage": "./Images/profile-picture-UID1.jpg",
      "role": "Student",
      "course": "Computer Science",
      "year": "3rd"
    }
  }
  ```
- **Example**:
  ```bash
  curl -X PUT http://localhost:3000/api/posts/1/update \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Post","description":"Updated description"}'
  ```

#### DELETE /api/posts/:id/delete
Delete a post.

- **Path Parameters**:
  - `id`: Post ID (integer)
- **Response (200 OK)**:
  ```json
  {"message": "Post deleted successfully"}
  ```
- **Example**:
  ```bash
  curl -X DELETE http://localhost:3000/api/posts/1/delete
  ```

### Job Endpoints
Manage job listings for campus opportunities.

#### GET /api/jobs
List all jobs.

- **Response (200 OK)**:
  ```json
  [
    {
      "id": 1,
      "createdAt": "2025-04-24T10:00:00Z",
      "updatedAt": "2025-04-24T10:00:00Z",
      "deletedAt": null,
      "title": "Software Engineer Intern",
      "description": "Internship at Tech Corp",
      "company": "Tech Corp",
      "link": "https://techcorp.com/jobs"
    },
    ...
  ]
  ```
- **Example**:
  ```bash
  curl http://localhost:3000/api/jobs
  ```

#### POST /api/jobs
Create a new job listing.

- **Request Body**:
  ```json
  {
    "title": "string (required)",
    "description": "string (required)",
    "company": "string (required)",
    "link": "string (required)"
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T10:00:00Z",
    "deletedAt": null,
    "title": "Software Engineer Intern",
    "description": "Internship at Tech Corp",
    "company": "Tech Corp",
    "link": "https://techcorp.com/jobs"
  }
  ```
- **Example**:
  ```bash
  curl -X POST http://localhost:3000/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"title":"Software Engineer Intern","description":"Internship at Tech Corp","company":"Tech Corp","link":"https://techcorp.com/jobs"}'
  ```

#### GET /api/jobs/:id
Get a job by ID.

- **Path Parameters**:
  - `id`: Job ID (integer)
- **Response (200 OK)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T10:00:00Z",
    "deletedAt": null,
    "title": "Software Engineer Intern",
    "description": "Internship at Tech Corp",
    "company": "Tech Corp",
    "link": "https://techcorp.com/jobs"
  }
  ```
- **Example**:
  ```bash
  curl http://localhost:3000/api/jobs/1
  ```

#### PUT /api/jobs/:id/update
Update a job listing.

- **Path Parameters**:
  - `id`: Job ID (integer)
- **Request Body**:
  ```json
  {
    "title": "string (optional)",
    "description": "string (optional)",
    "company": "string (optional)",
    "link": "string (optional)"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T12:00:00Z",
    "deletedAt": null,
    "title": "Senior Software Engineer Intern",
    "description": "Updated internship at Tech Corp",
    "company": "Tech Corp",
    "link": "https://techcorp.com/jobs"
  }
  ```
- **Example**:
  ```bash
  curl -X PUT http://localhost:3000/api/jobs/1/update \
  -H "Content-Type: application/json" \
  -d '{"title":"Senior Software Engineer Intern","description":"Updated internship at Tech Corp"}'
  ```

#### DELETE /api/jobs/:id/delete
Delete a job listing.

- **Path Parameters**:
  - `id`: Job ID (integer)
- **Response (200 OK)**:
  ```json
  {"message": "Job deleted successfully"}
  ```
- **Example**:
  ```bash
  curl -X DELETE http://localhost:3000/api/jobs/1/delete
  ```

### Event Endpoints
Manage campus events.

#### GET /api/events
List all events.

- **Response (200 OK)**:
  ```json
  [
    {
      "id": 1,
      "createdAt": "2025-04-24T10:00:00Z",
      "updatedAt": "2025-04-24T10:00:00Z",
      "deletedAt": null,
      "quarter": "Spring",
      "month": "April",
      "date": "2025-04-25T18:00:00Z",
      "title": "Career Fair",
      "participants": "All students"
    },
    ...
  ]
  ```
- **Example**:
  ```bash
  curl http://localhost:3000/api/events
  ```

#### POST /api/events
Create a new event.

- **Request Body**:
  ```json
  {
    "quarter": "string (required)",
    "month": "string (required)",
    "date": "string (ISO 8601, required, e.g., 2025-04-25T18:00:00Z)",
    "title": "string (required)",
    "participants": "string (optional)"
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T10:00:00Z",
    "deletedAt": null,
    "quarter": "Spring",
    "month": "April",
    "date": "2025-04-25T18:00:00Z",
    "title": "Career Fair",
    "participants": "All students"
  }
  ```
- **Example**:
  ```bash
  curl -X POST http://localhost:3000/api/events \
  -H "Content-Type: application/json" \
  -d '{"quarter":"Spring","month":"April","date":"2025-04-25T18:00:00Z","title":"Career Fair","participants":"All students"}'
  ```

#### GET /api/events/:id
Get an event by ID.

- **Path Parameters**:
  - `id`: Event ID (integer)
- **Response (200 OK)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T10:00:00Z",
    "deletedAt": null,
    "quarter": "Spring",
    "month": "April",
    "date": "2025-04-25T18:00:00Z",
    "title": "Career Fair",
    "participants": "All students"
  }
  ```
- **Example**:
  ```bash
  curl http://localhost:3000/api/events/1
  ```

#### PUT /api/events/:id/update
Update an event.

- **Path Parameters**:
  - `id`: Event ID (integer)
- **Request Body**:
  ```json
  {
    "quarter": "string (optional)",
    "month": "string (optional)",
    "date": "string (ISO 8601, optional)",
    "title": "string (optional)",
    "participants": "string (optional)"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T12:00:00Z",
    "deletedAt": null,
    "quarter": "Spring",
    "month": "April",
    "date": "2025-04-26T18:00:00Z",
    "title": "Updated Career Fair",
    "participants": "All students"
  }
  ```
- **Example**:
  ```bash
  curl -X PUT http://localhost:3000/api/events/1/update \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Career Fair","date":"2025-04-26T18:00:00Z"}'
  ```

#### DELETE /api/events/:id/delete
Delete an event.

- **Path Parameters**:
  - `id`: Event ID (integer)
- **Response (200 OK)**:
  ```json
  {"message": "Event deleted successfully"}
  ```
- **Example**:
  ```bash
  curl -X DELETE http://localhost:3000/api/events/1/delete
  ```

### Timetable Endpoints
Manage class timetables.

#### GET /api/timetables
List all timetables.

- **Response (200 OK)**:
  ```json
  [
    {
      "id": 1,
      "createdAt": "2025-04-24T10:00:00Z",
      "updatedAt": "2025-04-24T10:00:00Z",
      "deletedAt": null,
      "day": "Monday",
      "subject": "Data Structures",
      "subjectCode": "CS201",
      "faculty": "Engineering",
      "room": "A101",
      "time": "2025-04-24T09:00:00Z",
      "instructor": "Dr. Smith"
    },
    ...
  ]
  ```
- **Example**:
  ```bash
  curl http://localhost:3000/api/timetables
  ```

#### POST /api/timetables
Create a new timetable entry.

- **Request Body**:
  ```json
  {
    "day": "string (required)",
    "subject": "string (required)",
    "subjectCode": "string (required)",
    "faculty": "string (required)",
    "room": "string (required)",
    "time": "string (ISO 8601, required, e.g., 2025-04-24T09:00:00Z)",
    "instructor": "string (required)"
  }
  ```
- **Response (201 Created)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T10:00:00Z",
    "deletedAt": null,
    "day": "Monday",
    "subject": "Data Structures",
    "subjectCode": "CS201",
    "faculty": "Engineering",
    "room": "A101",
    "time": "2025-04-24T09:00:00Z",
    "instructor": "Dr. Smith"
  }
  ```
- **Example**:
  ```bash
  curl -X POST http://localhost:3000/api/timetables \
  -H "Content-Type: application/json" \
  -d '{"day":"Monday","subject":"Data Structures","subjectCode":"CS201","faculty":"Engineering","room":"A101","time":"2025-04-24T09:00:00Z","instructor":"Dr. Smith"}'
  ```

#### GET /api/timetables/:id
Get a timetable entry by ID.

- **Path Parameters**:
  - `id`: Timetable ID (integer)
- **Response (200 OK)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T10:00:00Z",
    "deletedAt": null,
    "day": "Monday",
    "subject": "Data Structures",
    "subjectCode": "CS201",
    "faculty": "Engineering",
    "room": "A101",
    "time": "2025-04-24T09:00:00Z",
    "instructor": "Dr. Smith"
  }
  ```
- **Example**:
  ```bash
  curl http://localhost:3000/api/timetables/1
  ```

#### PUT /api/timetables/:id/update
Update a timetable entry.

- **Path Parameters**:
  - `id`: Timetable ID (integer)
- **Request Body**:
  ```json
  {
    "day": "string (optional)",
    "subject": "string (optional)",
    "subjectCode": "string (optional)",
    "faculty": "string (optional)",
    "room": "string (optional)",
    "time": "string (ISO 8601, optional)",
    "instructor": "string (optional)"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "id": 1,
    "createdAt": "2025-04-24T10:00:00Z",
    "updatedAt": "2025-04-24T12:00:00Z",
    "deletedAt": null,
    "day": "Monday",
    "subject": "Advanced Data Structures",
    "subjectCode": "CS201",
    "faculty": "Engineering",
    "room": "A102",
    "time": "2025-04-24T10:00:00Z",
    "instructor": "Dr. Smith"
  }
  ```
- **Example**:
  ```bash
  curl -X PUT http://localhost:3000/api/timetables/1/update \
  -H "Content-Type: application/json" \
  -d '{"subject":"Advanced Data Structures","room":"A102","time":"2025-04-24T10:00:00Z"}'
  ```

#### DELETE /api/timetables/:id/delete
Delete a timetable entry.

- **Path Parameters**:
  - `id`: Timetable ID (integer)
- **Response (200 OK)**:
  ```json
  {"message": "Timetable deleted successfully"}
  ```
- **Example**:
  ```bash
  curl -X DELETE http://localhost:3000/api/timetables/1/delete
  ```

## Notes
- **Image Handling**:
  - Post and user profile images are stored in the `./Images` folder and served at `/Images`.
  - Images are uploaded as base64-encoded strings and saved with filenames based on the post title (e.g., `post-my-post.jpg`) or user ID (for profiles).
  - Access images via `http://localhost:3000/Images/<filename>`.
- **Time Formats**:
  - Use ISO 8601 format for `date` and `time` fields (e.g., `2025-04-24T09:00:00Z`).
- **Validation**:
  - Required fields are enforced (e.g., `title` for posts, `email` for users).
  - Invalid `userID` in posts will return a 400 error.
- **Testing**:
  - Use tools like Postman or curl to test endpoints.
  - Example Postman collection can be provided upon request.

## Contact
For questions or issues, contact the backend team:
- **Email**: kevinmranda042@gmail.com
- **github**: kevinmranda

---

*Generated on April 24, 2025*