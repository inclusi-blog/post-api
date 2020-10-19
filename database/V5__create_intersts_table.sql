create table INTERESTS
(
    INTEREST_ID BIGINT auto_increment,
    NAME VARCHAR(15) not null,
    CREATED_AT TIMESTAMP default current_timestamp not null,
    UPDATED_AT TIMESTAMP default NULL null,
    DELETED_AT TIMESTAMP default NULL null,
    constraint INTERESTS_pk
        primary key (INTEREST_ID)
);

create unique index INTERESTS_NAME_uindex
    on INTERESTS (NAME);

