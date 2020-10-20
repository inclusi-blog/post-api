alter table DRAFTS
    add TAGLINE VARCHAR(100) null after TITLE_DATA;

alter table DRAFTS modify UPDATED_AT timestamp null after CREATED_AT;

