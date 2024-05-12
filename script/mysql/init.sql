create database if not exists webook;
use webook;

create table if not exists user (
    UserID int primary key auto_increment not null ,
    UserName varchar(20),
    Email varchar(20),
    Password varchar(20),
    CreateTime bigint,
    UpdateTime bigint
);

create table if not exists role (
    RoleID int primary key auto_increment not null ,
    RoleName varchar(20),
    RolePid int,
    CreateTime bigint,
    UpdateTime bigint
);

create table if not exists user_role(
    ID int primary key auto_increment not null ,
    UserID int,
    RoleID int
);

drop table if exists router;
create table if not exists router(
    RouterID int primary key auto_increment not null ,
    RouterName varchar(20),
    RouterPath varchar(20),
    RouterPid int,
    RouterMethod varchar(20),
    RouterState int,
    CreateTime bigint,
    UpdateTime bigint
);

create table if not exists role_router(
    ID int primary key auto_increment not null ,
    RoleID int,
    RouterID int
);