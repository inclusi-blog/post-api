DROP SCHEMA IF EXISTS POST;
CREATE DATABASE POST;
DROP USER IF EXISTS POST;
CREATE USER 'POST'@'%' IDENTIFIED BY 'gola';
GRANT ALL PRIVILEGES ON POST.* TO 'POST'@'%';
FLUSH PRIVILEGES;
