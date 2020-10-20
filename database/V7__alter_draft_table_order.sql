alter table DRAFTS modify CREATED_AT timestamp default CURRENT_TIMESTAMP null after DELETE_AT;

alter table DRAFTS modify TITLE_DATA json null after POST_DATA;

