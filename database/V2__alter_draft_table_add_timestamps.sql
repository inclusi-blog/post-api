alter table DRAFTS modify DRAFT_ID varchar(12) not null;

alter table DRAFTS
    add CREATED_AT timestamp default current_timestamp null;

alter table DRAFTS
    add UPDATED_AT timestamp null;

alter table DRAFTS
    add DELETE_AT timestamp null;

