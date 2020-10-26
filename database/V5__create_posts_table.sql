create table POSTS
(
    ID bigserial not null
        constraint posts_pk
        primary key,
    PUID varchar(12),
    USER_ID bigint not null,
    POST_DATA json not null,
    TITLE_DATA json not null,
    READ_TIME int not null,
    VIEW_COUNT int default 0 not null,
    CREATED_AT timestamp default current_timestamp not null,
    UPDATED_AT timestamp,
    DELETED_AT timestamp
);

