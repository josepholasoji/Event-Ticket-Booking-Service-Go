-- Active: 1722204061621@@127.0.0.1@3306@booking

-- create table event
CREATE TABLE IF NOT EXISTS events(  
    id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    name VARCHAR(100) NOT NULL COMMENT 'Name of the event',
    description TEXT NOT NULL COMMENT 'Description of the event',
    capacity INT NOT NULL COMMENT 'Capacity of the event',
    start_date DATE NOT NULL COMMENT 'Start date of the event',
    category ENUM('Concert', 'Conference', 'Game') NOT NULL COMMENT 'Category of the event',
    end_date DATE NOT NULL COMMENT 'End date of the event',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at'
) COMMENT 'Event Table';

-- create table user
CREATE TABLE IF NOT EXISTS users(  
    id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    name VARCHAR(100) NOT NULL COMMENT 'Name of the user',
    email VARCHAR(100) NOT NULL COMMENT 'Email of the user',
    password VARCHAR(200) NOT NULL COMMENT 'Password of the user',
    role ENUM('admin', 'user') DEFAULT 'user' COMMENT 'Role of the user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated at',
    is_active BOOLEAN DEFAULT TRUE COMMENT 'Is the user active?'
) COMMENT 'User Table';

-- create a table of user event booking
CREATE TABLE IF NOT EXISTS bookings(
    id int NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    user_id int NOT NULL COMMENT 'User ID',
    event_id int NOT NULL COMMENT 'Event ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Created at'
) COMMENT 'User Event Booking Table';

-- Event notfication table
CREATE TABLE IF NOT EXISTS event_notifications (
    id SERIAL PRIMARY KEY,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);