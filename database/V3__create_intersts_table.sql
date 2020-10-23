create table INTERESTS
(
    ID bigserial not null
        constraint interests_pk
        primary key,
    NAME varchar(15) not null,
    CREATED_AT timestamp default current_timestamp not null,
    UPDATED_AT timestamp,
    DELETED_AT timestamp
);

create unique index interests_name_uindex
    on INTERESTS (NAME);

