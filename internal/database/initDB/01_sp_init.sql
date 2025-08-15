DELIMITER $$

-- Get's
CREATE PROCEDURE sp_get_posts_by_ids_direct(IN ids_csv TEXT)
BEGIN
    IF ids_csv NOT REGEXP '^[0-9]+(,[0-9]+)*$' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid IDs input';
    END IF;

    SET @sql_stmt = CONCAT(
        'SELECT id, sender_id, content, created_date
         FROM posts
         WHERE id IN (', ids_csv, ')'
    );
    PREPARE stmt FROM @sql_stmt;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    SET @sql_stmt = NULL;
END$$

CREATE PROCEDURE sp_get_posts_by_ids_temp(IN ids_csv TEXT)
BEGIN
    IF ids_csv NOT REGEXP '^[0-9]+(,[0-9]+)*$' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid IDs input';
    END IF;

    CREATE TEMPORARY TABLE tmp_ids (
        id INT PRIMARY KEY
    );

    SET @insert_sql = CONCAT(
        'INSERT INTO tmp_ids (id) VALUES ',
        REPLACE(ids_csv, ',', '),(')
    );

    PREPARE stmt FROM @insert_sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    SELECT p.id, p.sender_id, p.content, p.created_date
    FROM posts p
    JOIN tmp_ids t ON p.id = t.id;

    DROP TEMPORARY TABLE IF EXISTS tmp_ids;

    SET @insert_sql = NULL;
END$$

-- Insert's :
CREATE PROCEDURE sp_insert_posts_bulk(IN posts_json JSON)
BEGIN
    INSERT INTO posts (sender_id, content, created_date)
    SELECT sender_id, content, NOW()
    FROM JSON_TABLE(
        posts_json,
        "$[*]" COLUMNS (
            sender_id INT PATH "$.sender_id",
            content  TEXT PATH "$.content"
        )
    ) AS jt;
END$$

DELIMITER ;
