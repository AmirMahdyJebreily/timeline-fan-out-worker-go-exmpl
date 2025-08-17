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

CREATE PROCEDURE sp_get_subscribers_by_user(IN p_user_id INT)
BEGIN
    SELECT subscriber_id
    FROM subscriber_users
    WHERE sender_id = p_user_id;
END$$

CREATE PROCEDURE sp_get_subscribers_bulk(IN p_sender_ids_json JSON)
BEGIN
  IF p_sender_ids_json IS NULL OR JSON_LENGTH(p_sender_ids_json) = 0 THEN
    SELECT CAST(NULL AS UNSIGNED) AS sender_id, CAST(NULL AS UNSIGNED) AS subscriber_id
      WHERE FALSE;
    RETURN;
  END IF;

  SELECT s.sender_id,
         s.subscriber_id
  FROM subscriber_users AS s
  JOIN (
    SELECT DISTINCT CAST(jt.sender_id AS UNSIGNED) AS sender_id
    FROM JSON_TABLE(
      p_sender_ids_json,
      '$[*]' COLUMNS (
        sender_id BIGINT PATH '$'
      )
    ) AS jt
    WHERE jt.sender_id IS NOT NULL
  ) AS ul ON s.sender_id = ul.sender_id
  ORDER BY s.sender_id, s.subscriber_id;
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
