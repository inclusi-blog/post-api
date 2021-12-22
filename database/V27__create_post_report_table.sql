create table post_reports
(
    post_id     uuid                                  not null
        constraint post_reports_posts_id_fk
            references posts,
    reported_by uuid                                  not null
        constraint post_reports_users_id_fk
            references users,
    reason      text,
    created_at  timestamptz default current_timestamp not null,
    constraint post_reports_pk
        primary key (post_id, reported_by)
);

create unique index post_reports_post_id_reported_by_uindex
    on post_reports (post_id, reported_by);

