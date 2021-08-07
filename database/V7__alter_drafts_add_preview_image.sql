alter table drafts drop column created_at;

alter table drafts
    add preview_image text;

alter table drafts
    add created_at timestamp default current_timestamp not null;

alter table drafts drop column updated_at;

alter table drafts drop column deleted_at;

alter table drafts
    add updated_at timestamp;

alter table drafts
    add deleted_at timestamp;

