## 4.1 Setting up MySQL
The official MySQL documentation contains comprehensive [installation instructions](https://dev.mysql.com/doc/refman/8.0/en/installing.htmls) for all types of operating systems, but if you’re using Mac OS you should be able to install it with:
```
$ brew install mysql
```

Or if you’re using a Linux distribution which supports apt (like Debian and Ubuntu) you can install it with:
```
$ sudo apt install mysql-server
```

### Scaffolding the database
Once connected, the 1st thing we need to do is establish a database in MySQL to store all the data for our project. Copy and paste the following commands into the mysql prompt to create a new snippetbox database using UTF8 encoding.
```sql
-- Create a new UTF-8 `snippetbox` database.
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Switch to using the `snippetbox` database.
USE snippetbox;
```

Then copy and paste the following SQL statement to create a new snippets table to hold the text snippets for our application:
```sql
-- Create a `snippets` table.
CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL
);

-- Add an index on the created_at column.
CREATE INDEX idx_snippets_created_at ON snippets(created_at);
```

Each record in this table will have an integer `id` field which will act as the unique identifier for the text snippet. It will also have a short text `title` and the snippet content itself will be stored in the `content` field. We’ll also keep some metadata about the times that the snippet was `created_at` and when it `expires_at`.

```sql
-- Add some dummy records.
INSERT INTO snippets (title, content, created_at, expires_at) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title, content, created_at, expires_at) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title, content, created_at, expires_at) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);
```

### Creating a new user
From a security point of view, it’s not a good idea to connect to MySQL as the root user from a web application. Instead it’s better to create a database user with restricted permissions on the database.

So, while you’re still connected to the MySQL prompt run the following commands to create a new web user with SELECT, INSERT, UPDATE and DELETE privileges only on the database.

```sql
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
-- Important: Make sure to swap 'pass' with a password of your own choosing.
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';
```

### Test the new user
You should now be able to connect to the snippetbox database as the web user using the following command. When prompted enter the password that you just set.
```
$ mysql -D snippetbox -u web -p
Enter password:
mysql>
```

If the permissions are working correctly you should find that you’re able to perform SELECT and INSERT operations on the database correctly, but other commands such as DROP TABLE and GRANT will fail.

```
mysql> SELECT id, title, expires_at FROM snippets;
+----+------------------------+---------------------+
| id | title                  | expires_at          |
+----+------------------------+---------------------+
|  1 | An old silent pond     | 2025-03-18 10:00:26 |
|  2 | Over the wintry forest | 2025-03-18 10:00:26 |
|  3 | First autumn morning   | 2024-03-25 10:00:26 |
+----+------------------------+---------------------+
3 rows in set (0.00 sec)
```
```
mysql> DROP TABLE snippets;
ERROR 1142 (42000): DROP command denied to user 'web'@'localhost' for table 'snippets'
```

---
## 4.2 Installing a database driver
To use MySQL from our Go web application we need to install a `database driver`. This essentially acts as a middleman, translating commands between Go and the MySQL database itself.

You can find a comprehensive [list of available drivers](https://go.dev/wiki/SQLDrivers) on the Go wiki, but for our application we’ll use the popular [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) driver.

