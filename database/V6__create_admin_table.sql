create table admin
(
    id uuid not null
        constraint admin_pk
        primary key,
    name varchar(50) not null,
    email varchar(60) not null,
    role_id uuid not null
        constraint admin_roles_id_fk
        references roles
);

