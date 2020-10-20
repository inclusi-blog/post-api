create table POST_INTERESTS
(
    ID BIGINT auto_increment,
    POST_ID VARCHAR(12) not null,
    INTEREST BIGINT not null,
    constraint POST_INTERESTS_pk
        primary key (ID),
    constraint POST_INTERESTS_INTERESTS_INTEREST_ID_fk
        foreign key (INTEREST) references INTERESTS (INTEREST_ID)
);

