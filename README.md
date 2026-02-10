**Event Booking App:**

A Go + Buffalo–based microservice event booking and notification platform that enables users to register, authenticate, create and manage events, search events, book tickets, and receive event notifications. The system is built using Go (1.21+), Buffalo framework, Docker Compose, MariaDB, and a message broker (Kafka), and exposes RESTful JSON APIs for user management, authentication, event management, booking workflows, and notification delivery.

The repo includes unit and integration tests, containerized local development setup, and seeded sample data for rapid onboarding and testing.

FUNTIONAL REQUIREMENTS:

1.  Create an account;
2.  User authentication to log into the system;
3.  Create events;
4.  Search and reserve tickets for events;
5.  Send notification to users before event starts.
6.  Event users should subscribe, listen and get notification messages.

**_NB:_** _To test out the implemented functional requitements: Simply run the “_**_functionalRequirementTest_**_” integration test inside the \`_**_UserApiRestControllerTests_**_.java\` file to see all the implemented functional requirements running._

NON-FUNCTIONAL REQUIREMENTS:

1.  The project MUST be buildable and runnable;
2.  The project MUST have Unit tests;
3.  The project MUST have a README file with build/run/test instructions (use a DB that can be run locally, e.g. in-memory, via container);
4.  Any data required by the application to run (e.g. reference tables, dummy data) MUST be preloaded in the database;
5.  Input/output data MUST be in JSON format;
6.  Use a framework of your choice, but popular, up-to-date, and long-term support versions are recommended.

## REQUIREMENTS TO RUN THE APPLICATION

1. **Go (v1.21 or later)**  
   Required to build and run the Buffalo-based services.

2. **Buffalo CLI**  
   Required for running and managing the Buffalo application.

3. **Docker Desktop**  
   Required for running supporting services (e.g., database, message broker).
   - **Windows users:** Ensure **WSL 2** is enabled and configured.

---

**Note:**  
Ensure `go`, `buffalo`, `docker`, and `docker compose` are available in your system PATH before starting the application.

## RUNNING AND TESTING THE APPLICATION

1. **Navigate to the service directory**

   ```bash
   cd <base-folder>/bookingservice

   ```

2. **Start dependent services using Docker Compose**
   ```bash
   docker compose -f compose.yaml up -d
   ```
This will start required supporting services:

- 1. MariaDB — Stores application data
- 2. Message Broker (e.g., Kafka) — Handles event notifications

3. **Run tests and build the service**
   ```bash
   go test ./...
   go build -o bookingservice
   ```
- This step runs all tests and prepares the service binary (and can trigger sample data seeding if configured).

4. **Run the service**
    ````bash
    ./bookingservice
5. **To build, run amd auto reload while developing you service**
    ```bash
    buffalo dev
    ```
ASSUMTIONS MADE IN THIS DEMO APP PROJECT:
1. This is a test project therefore:
2.  There is no requirement for an accurate event timing and notification.
3.  There is no requirement to test for a in ability to book an even due to capacity

ANY ISSUES:
Please, ensure that other app is not using port 8080.
API ENDPOINTS: _See the added pictures for screenshot of local postman calls._
Base-path: https://localhost:8080/booking/api

|     |     |     |
| --- | --- | --- |
| S/N | Endpoint | Description |
| 1   | POST /&lt;base-path&gt;/users | Create a user |
| 2   | GET /&lt;base-path&gt;/users/events | Get a user booked events |
| 3   | POST /&lt;base-path&gt;/events | Create an event |
| 4   | GET POST /&lt;base-path&gt;/events/{searchPhrase} | Search created events |
| 5   | GET /&lt;base-path&gt;/events/{eventId} | Get event by Id |
| 6   | GET POST /&lt;base-path&gt;/events/{eventId}/tickets | Book an event |
| 7   | DELETE /&lt;base-path&gt;/events/{eventid}/tickets/{ticketsId} | Cancel a booking |
| 8   | POST /&lt;base-path&gt;/auth | Login a user |

DEMO (SAMPLE) PRELOADED USER:
Username: sample@admin.com
Password: adminadmin

DATA AND TABLE DEFINITIONS:
Please see, init-scipts for the SQL migration files.
<br/>bookings

|     |     |     |     |     |     |
| --- | --- | --- | --- | --- | --- |
| **Column** | **Type** | **Comment** | **PK** | **Nullable** | **Default** |
| id  | int(11) | Primary Key | YES | NO  |     |
| user_id | int(11) | User ID |     | NO  |     |
| event_id | int(11) | Event ID |     | NO  |     |
| created_at | timestamp | Created at |     | YES | current_timestamp() |
| updated_at | timestamp | Updated at |     | YES | current_timestamp() |

events

|     |     |     |     |     |     |
| --- | --- | --- | --- | --- | --- |
| **Column** | **Type** | **Comment** | **PK** | **Nullable** | **Default** |
| id  | int(11) | Primary Key | YES | NO  |     |
| name | varchar(255) |     |     | YES | NULL |
| description | varchar(255) |     |     | YES | NULL |
| capacity | int(11) | Capacity of the event |     | NO  |     |
| start_date | date | Start date of the event |     | NO  |     |
| category | tinyint(4) |     |     | YES | NULL |
| end_date | datetime(6) |     |     | YES | NULL |
| created_at | timestamp | Created at |     | YES | current_timestamp() |
| updated_at | timestamp | Updated at |     | YES | current_timestamp() |
| date | datetime(6) |     |     | YES | NULL |
| is_active | bit(1) |     |     | YES | NULL |

event_notifications

|     |     |     |     |     |     |
| --- | --- | --- | --- | --- | --- |
| **Column** | **Type** | **Comment** | **PK** | **Nullable** | **Default** |
| id  | bigint(20) unsigned |     | YES | NO  |     |
| event_id | int(11) |     |     | NO  |     |
| user_id | int(11) |     |     | NO  |     |
| created_at | timestamp |     |     | NO  | current_timestamp() |
| updated_at | timestamp |     |     | NO  | current_timestamp() |

users

|     |     |     |     |     |     |
| --- | --- | --- | --- | --- | --- |
| **Column** | **Type** | **Comment** | **PK** | **Nullable** | **Default** |
| id  | int(11) | Primary Key | YES | NO  |     |
| name | varchar(255) |     |     | YES | NULL |
| email | varchar(255) |     |     | YES | NULL |
| password | varchar(255) |     |     | YES | NULL |
| role | varchar(255) |     |     | YES | NULL |
| created_at | timestamp | Created at |     | YES | current_timestamp() |
| updated_at | timestamp | Updated at |     | YES | current_timestamp() |
| is_active | tinyint(1) | Is the user active? |     | YES | 1   |
````
