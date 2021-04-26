CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password  TEXT NOT NULL,
);

CREATE TABLE IF NOT EXISTS black_list(
    token TEXT NOT NULL UNIQUE
);

INSERT INTO users(username, password)
    VALUES
        ('cfabrica46',  '01234'),
        ('arturo',      '12345'),
        ('carlos',      'abcd'),
        ('luis',        'lolsito123');