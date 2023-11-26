create table draft_images
(
    id uuid not null,
    draft_id uuid not null,
    upload_id uuid not null
);

create unique index draft_images_upload_id_draft_id_uindex
    on draft_images (draft_id, upload_id);

create unique index draft_images_id_uindex
    on draft_images (id);

alter table draft_images
    add constraint draft_images_pk
        primary key (id);

alter table draft_images
    add constraint draft_images_drafts_id_fk
        foreign key (draft_id) references drafts;


