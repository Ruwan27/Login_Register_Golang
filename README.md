# Login_Register_Golang
use below query create database 
CREATE DATABASE login_example;

USE login_example;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

#Go run command | use one of them 
go run ./main.go 
go run ./
go run .
