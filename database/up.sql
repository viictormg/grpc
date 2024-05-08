DROP TABLE IF EXISTS students;

CREATE TABLE students (
    id VARCHAR(32) PRIMARY KEY,
    name varchar(255) NOT NULL,
    age INTEGER NOT NULL
);

DROP TABLE IF EXISTS tests;

CREATE TABLE tests (
    id VARCHAR(32) PRIMARY KEY,
    name varchar(255) NOT NULL
);

DROP TABLE IF EXISTS questions; 

CREATE TABLE questions (
    id VARCHAR(32) PRIMARY KEY,
    test_id varchar(255) NOT NULL,
    question varchar(255) NOT NULL,
    answer varchar(255) NOT NULL,
    CONSTRAINT fk_test FOREIGN KEY(test_id) REFERENCES tests(id)
)