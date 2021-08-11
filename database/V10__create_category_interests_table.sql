create table category_x_interests
(
    category_id uuid
        constraint category_x_interests_interests_id_fk
        references interests,
    interest_id uuid
        constraint category_x_interests_interests_id_fk_2
        references interests,
    constraint category_x_interests_pk
        primary key (category_id, interest_id)
);

