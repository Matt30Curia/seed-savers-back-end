CREATE TABLE IF NOT EXISTS users (
    user_id INT AUTO_INCREMENT NOT NULL UNIQUE,
    name VARCHAR(50),
    email VARCHAR(100)  NOT NULL,
    password VARCHAR(100) ,
    credits SMALLINT DEFAULT 1,
    PRIMARY KEY (user_id )
    
);