create table DRAFTS
(
    ID BIGINT auto_increment,
    DRAFT_ID VARCHAR(10) not null,
    USER_ID BIGINT not null,
    POST_DATA JSON null,
    constraint DRAFTS_pk
        primary key (ID)
);

create unique index DRAFTS_DRAFT_ID_uindex
    on DRAFTS (DRAFT_ID);

create unique index DRAFTS_USER_ID_uindex
    on DRAFTS (USER_ID);

