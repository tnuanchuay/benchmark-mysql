package main

const CreateWriteTable = `
CREATE TABLE IF NOT EXISTS benchmark_write (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ip VARCHAR(50) NOT NULL,
    created_at TIMESTAMP
);
`

const CreateReadTable = `
CREATE TABLE IF NOT EXISTS benchmark_read (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ip VARCHAR(50) NOT NULL,
    url VARCHAR(50) NOT NULL,
    redirect_to VARCHAR(50) NOT NULL,
    created_at TIMESTAMP
);
`

const CreateSimpleData = `
INSERT INTO benchmark_read(ip, url, redirect_to, created_at) VALUES (?, ?, ?, ?)
`
