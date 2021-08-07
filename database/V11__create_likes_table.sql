create table LIKES
(
    ID bigserial not null
        constraint likes_pk
        primary key,
    LIKED_BY TEXT[],
    POST_ID bigserial,
    CONSTRAINT fk_likes
      FOREIGN KEY(POST_ID) 
	  REFERENCES POSTS(ID),
    CREATED_AT timestamp default current_timestamp not null,
    UPDATED_AT timestamp,
    DELETED_AT timestamp
);