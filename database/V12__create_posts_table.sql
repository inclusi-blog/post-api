create table posts
(
    id uuid,
    data jsonb not null,
    is_series boolean default false not null,
    is_under_publication boolean default false not null,
    author_id uuid not null,
    draft_id uuid not null,
    created_at timestamptz default current_timestamp not null,
    updated_at timestamptz,
    deleted_at timestamptz
);

create unique index posts_author_id_draft_id_uindex
    on posts (author_id, draft_id);

create unique index posts_id_uindex
    on posts (id);

alter table posts
    add constraint posts_pk
        primary key (id);

alter table posts
    add constraint posts_drafts_id_fk
        foreign key (draft_id) references drafts;

alter table posts
    add constraint posts_users_id_fk
        foreign key (author_id) references users;


