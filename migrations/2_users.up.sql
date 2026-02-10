

-- a procedure to add a new user into the user table
-- Parameters: name, email, password, role, created_at, updated_at, is_active

-- Check if the stored procedure exists
DELIMITER $$
DROP PROCEDURE IF EXISTS add_user;
CREATE PROCEDURE add_user(
    IN p_name VARCHAR(100),
    IN p_email VARCHAR(100),
    IN p_password VARCHAR(200),
    IN p_role ENUM('admin', 'user'),
    OUT p_id INT
)

BEGIN
    INSERT INTO users(name, email, password, role, created_at, updated_at, is_active)
    VALUES(p_name, p_email, p_password, p_role, NOW(), NOW(), TRUE);
    SET p_id = LAST_INSERT_ID();
END$$

DELIMITER ;

-- a procedure to get all users from the user table
-- Parameters: None
DELIMITER $$
DROP PROCEDURE IF EXISTS get_users;
CREATE PROCEDURE get_users()
BEGIN
    SELECT * FROM users;
END$$

DELIMITER ;

-- a procedure to get a user by id from the user table
-- Parameters: id
DELIMITER $$
DROP PROCEDURE IF EXISTS get_user_by_id;
CREATE PROCEDURE get_user_by_id(
    IN p_id INT
)

BEGIN
    SELECT * FROM users WHERE id = p_id;
END$$

DELIMITER ;

-- a procedure to delete a user by id from the user table
-- Parameters: id
DELIMITER $$
DROP PROCEDURE IF EXISTS delete_user_by_id;
CREATE PROCEDURE delete_user_by_id(
    IN p_id INT
)

BEGIN
    DELETE FROM bookings WHERE user_id = p_id;
    DELETE FROM event_notifications WHERE user_id = p_id;
    DELETE FROM users WHERE id = p_id;
END$$

DELIMITER ;

-- a procedure to search for a user by email from the user table
-- Parameters: email
DELIMITER $$
DROP PROCEDURE IF EXISTS get_user_by_email;
CREATE PROCEDURE get_user_by_email(
    IN p_email VARCHAR(100)
)

BEGIN
    SELECT * FROM users WHERE email = p_email;
END$$

DELIMITER ;

-- a procedure to get a user by email and password from the user table
-- Parameters: email, password
DELIMITER $$
DROP PROCEDURE IF EXISTS get_user_by_email_and_password;
CREATE PROCEDURE get_user_by_email_and_password(
    IN p_email VARCHAR(100),
    IN p_password VARCHAR(200)
)

BEGIN
    SELECT * FROM users WHERE email = p_email AND password = p_password;
END$$

DELIMITER ;

-- a procedure to delete all users from the user table
-- Parameters: None
DELIMITER $$
DROP PROCEDURE IF EXISTS delete_all_users;
CREATE PROCEDURE delete_all_users()
BEGIN
    DELETE FROM users;
END$$

DELIMITER ;

DELIMITER ;

-- a procedure to get the registered user events by id
-- Parameters: id
DELIMITER $$
DROP PROCEDURE IF EXISTS get_user_events_by_id;
CREATE PROCEDURE get_user_events_by_id(
    IN p_id INT
)

BEGIN
    SELECT * FROM events WHERE user_id = p_id;
END$$

DELIMITER ;