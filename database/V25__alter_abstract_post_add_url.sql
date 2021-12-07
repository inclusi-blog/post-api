alter table abstract_post
    add url text;

create unique index abstract_post_url_uindex
    on abstract_post (url);

