create database RedBook;
use redbook;
create table user (
    id bigint,
    email varchar(20),
    password varchar(20),
    createTime bigint,
    updateTime bigint
);