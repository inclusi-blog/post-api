insert into users(id, email, role_id) values (uuid_generate_v4(), 'hariharan@mensuvadi.com', (select id from roles where name = 'Admin'));
