**Event Booking App:**

A Spring Boot–based event booking and notification system that allows users to register, authenticate, create events, search and book tickets, and receive event notifications. The application is built with Java 17, Maven, Docker Compose, MariaDB, and ActiveMQ, includes unit and integration tests, preloaded sample data, and exposes RESTful JSON APIs for user, event, booking, and authentication workflows.

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

REQUIREMENT TO RUN THE APP:  
1\. JDK v17 (LTS). The sample is built with Spring Boot 3.xxx

2\. Maven 3x

3\. Docker desktop (WSL 2.0 for windows OS)

NB: Ensure Java, mvn and docker-compose are in your system path before trying to start the demo app.

TESTING THE APP:  
1\. CD into the **&lt;base folder&gt;/bookings**

**2\. Run: docker-compose -f .\\compose.yaml up -d**

This demo app uses the docker compose to load the dependent services:  
1\. Maria DB: To store the booking app database

2\. Active MQ: For event notification

**3\. Run: mvn clean;mvn package**

We need this step to run all test and generate sample data. Take note that the tests also contain the “_functionalRequirementTest_” which does run integrations test to confirm all the implemented functional requirements.

**4\. CD into &lt;base folder&gt;/bookings/target**

**5\. Run: java -jar booking-0.0.1-SNAPSHOT.jar**

ASSUMTIONS MADE IN THIS DEMO APP PROJECT:  
1\. This is a test project therefore:

1.  There is no requirement for an accurate event timing and notification.
2.  There is no requirement to test for a in ability to book an even due to capacity

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