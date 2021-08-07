create table DRAFTS
(
    ID bigserial not null
        constraint drafts_pk
        primary key,
    DRAFT_ID varchar(12) not null,
    USER_ID bigint not null,
    POST_DATA json,
    TITLE_DATA json,
    TAGLINE varchar(100),
    INTEREST json,
    CREATED_AT timestamp default current_timestamp not null,
    UPDATED_AT timestamp,
    DELETED_AT timestamp
);

create unique index drafts_draft_id_uindex
    on DRAFTS (DRAFT_ID);

