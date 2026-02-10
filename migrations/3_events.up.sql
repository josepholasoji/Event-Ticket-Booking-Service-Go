-- Active: 1722204061621@@127.0.0.1@3306@booking
-- a procedure to add an event
DELIMITER $$
DROP PROCEDURE IF EXISTS add_event;
CREATE PROCEDURE add_event(
    IN p_name VARCHAR(100),
    IN p_description TEXT,
    IN p_capacity INT,
    IN p_start_date DATE,
    IN p_end_date DATE,
    IN p_category ENUM('Concert', 'Conference', 'Game'),
    OUT p_id INT
)

BEGIN
    INSERT INTO events(name, description, capacity, start_date, end_date, category, created_at, updated_at, is_active)
    VALUES(p_name, p_description, p_capacity, p_start_date, p_end_date, p_category, NOW(), NOW(), 1);
    SET p_id = LAST_INSERT_ID();
END$$

DELIMITER ;

-- a procedure to get all events
DELIMITER $$
DROP PROCEDURE IF EXISTS get_events;
CREATE PROCEDURE get_events()
BEGIN
    SELECT * FROM events;
END$$

DELIMITER ;

-- a procedure to get an event by id

DELIMITER $$
DROP PROCEDURE IF EXISTS get_event_by_id;
CREATE PROCEDURE get_event_by_id(
    IN p_id INT
)

BEGIN
    SELECT * FROM events WHERE id = p_id;
END$$

DELIMITER ;

-- a procedure to delete an event by id

DELIMITER $$
DROP PROCEDURE IF EXISTS delete_event_by_id;
CREATE PROCEDURE delete_event_by_id(
    IN p_id INT
)

BEGIN
    DELETE FROM events WHERE id = p_id;
END$$

DELIMITER ;

-- a procedure to get an event by name

DELIMITER $$
DROP PROCEDURE IF EXISTS get_event_by_name;

CREATE PROCEDURE get_event_by_name(
    IN p_name VARCHAR(100)
)

BEGIN
    SELECT * FROM events WHERE name = p_name;
END$$

DELIMITER ;

-- a procedure to search for an event by name

DELIMITER $$
DROP PROCEDURE IF EXISTS search_event_by_name;
CREATE PROCEDURE search_event_by_name(
    IN p_name VARCHAR(100)
)

BEGIN
    SELECT * FROM events WHERE name = p_name;
END$$

DELIMITER ; 

-- a procedure to search for an event by partial name

DELIMITER $$
DROP PROCEDURE IF EXISTS search_event_by_partial_name;
CREATE PROCEDURE search_event_by_partial_name(
    IN p_name VARCHAR(100)
)

BEGIN
    SELECT * FROM events WHERE name LIKE CONCAT('%', p_name, '%');
END$$

DELIMITER ;

-- a precedure to delete all events

DELIMITER $$
DROP PROCEDURE IF EXISTS delete_all_events;
CREATE PROCEDURE delete_all_events()
BEGIN
    DELETE FROM events;
END$$

DELIMITER ;

-- a procedure to delete an event by id

DELIMITER $$
DROP PROCEDURE IF EXISTS delete_event_by_id;
CREATE PROCEDURE delete_event_by_id(
    IN p_id INT
)

BEGIN
    DELETE FROM events WHERE id = p_id;
END$$

DELIMITER ;


-- a stored procedure to book an event

DELIMITER $$

DROP PROCEDURE IF EXISTS book_event;

CREATE PROCEDURE book_event(
    IN p_event_id INT,
    IN p_user_id INT,
    OUT p_id INT
)

BEGIN
    DECLARE localId INT;
    SELECT id INTO localId FROM bookings WHERE event_id = p_event_id AND user_id = p_user_id;
    -- mariadb syntax to return the id of the booking if it already exists
    IF localId IS NOT NULL THEN
        SET p_id = localId;
    ELSE
        INSERT INTO bookings(event_id, user_id, created_at)
        VALUES(p_event_id, p_user_id, NOW());
        SET p_id = LAST_INSERT_ID();
    END IF;
END$$

DELIMITER ;

-- a stored procedure to get all bookings by user id

DELIMITER $$

DROP PROCEDURE IF EXISTS get_bookings_by_user_id;

CREATE PROCEDURE get_bookings_by_user_id(
    IN p_user_id INT
)

BEGIN
    SELECT * FROM bookings WHERE user_id = p_user_id;
END$$

DELIMITER ;

-- a stored procedure to get all bookings by event id

DELIMITER $$

DROP PROCEDURE IF EXISTS get_bookings_by_event_id;

CREATE PROCEDURE get_bookings_by_event_id(
    IN p_event_id INT
)

BEGIN
    SELECT * FROM bookings WHERE event_id = p_event_id;
END$$

DELIMITER ;

-- a stored procedure to get all events

DELIMITER $$
DROP PROCEDURE IF EXISTS get_all_events;

CREATE PROCEDURE get_all_events()

BEGIN
    SELECT * FROM events;
END$$

DELIMITER ;

-- a stored procedure to delete a booking by event id and user id

DELIMITER $$
DROP PROCEDURE IF EXISTS delete_booking_by_event_id_and_ticket_id;

CREATE PROCEDURE delete_booking_by_event_id_and_ticket_id(
    IN p_event_id INT,
    IN p_ticket_id INT
)

BEGIN
    DELETE FROM bookings WHERE event_id = p_event_id AND id = p_ticket_id;
END$$

DELIMITER ;

-- a stored procedure to get a booking by event id and user id

DELIMITER $$
DROP PROCEDURE IF EXISTS get_booking_by_event_id_and_user_id;

CREATE PROCEDURE get_booking_by_event_id_and_user_id(
    IN p_event_id INT,
    IN p_user_id INT,
    OUT p_id INT
)

BEGIN
    SELECT id INTO p_id FROM bookings WHERE event_id = p_event_id AND user_id = p_user_id;
END$$

DELIMITER ;

-- a stored procedure to get if an event has started

DELIMITER $$
DROP PROCEDURE IF EXISTS is_event_started;

CREATE PROCEDURE is_event_started(
    IN p_event_id INT,
    OUT p_has_started BOOLEAN
)

BEGIN
    DECLARE localStartDate DATE;
    SELECT start_date INTO localStartDate FROM events WHERE id = p_event_id;
    IF localStartDate < NOW() THEN
        SET p_has_started = TRUE;
    ELSE
        SET p_has_started = FALSE;
    END IF;
END$$

DELIMITER ;

-- a stored procedure to get if an event is almost starting

DELIMITER $$

DROP PROCEDURE IF EXISTS is_event_about_to_start;

CREATE PROCEDURE is_event_about_to_start(
    IN p_event_id INT,
    OUT p_is_almost_starting BOOLEAN
)

BEGIN
    DECLARE localStartDate DATE;
    SELECT start_date INTO localStartDate FROM events WHERE id = p_event_id;
    IF localStartDate < NOW() + INTERVAL 10 MINUTE THEN
        SET p_is_almost_starting = TRUE;
    ELSE
        SET p_is_almost_starting = FALSE;
    END IF;
END$$

DELIMITER ;

-- a stored procedure to log an event notification using the event id and the user id, we firstly get the number of event notifications for the user and the event, if it is greater than p_limit then we return p_limit and skip the insertion, otherwise we insert the event notification

DELIMITER $$
DROP PROCEDURE IF EXISTS log_event_notification;

CREATE PROCEDURE log_event_notification(
    IN p_event_id INT,
    IN p_user_id INT,
    IN p_limit INT,
    OUT p_count INT
)

BEGIN
    DECLARE localCount INT;
    SELECT COUNT(*) INTO localCount FROM event_notifications WHERE event_id = p_event_id AND user_id = p_user_id;
    IF localCount >= p_limit THEN
        SET p_count = p_limit;
    ELSE
        INSERT INTO event_notifications(event_id, user_id, created_at, updated_at)
        VALUES(p_event_id, p_user_id, NOW(), NOW());
        SET p_count = localCount + 1;
    END IF;
END$$

DELIMITER ;