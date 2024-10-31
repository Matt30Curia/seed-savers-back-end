CREATE TABLE IF NOT EXISTS adress(
    id INT NOT NULL PRIMARY KEY,
    state VARCHAR(40) NOT NULL,
    city VARCHAR(40) NOT NULL,
    street VARCHAR(40) NOT NULL,
    cap VARCHAR(10) NOT NULL,
    province VARCHAR(40) NOT NULL,
    number SMALLINT NOT NULL,
    apartment_number  VARCHAR(40) ,
    FOREIGN KEY (id ) REFERENCES users(user_id) ON DELETE CASCADE
);