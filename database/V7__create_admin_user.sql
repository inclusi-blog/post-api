insert into admin(id, name, email, role_id) values (uuid_generate_v4(), 'Hariharan', 'hariharan@mensuvadi.com', (select id from roles where name = 'Admin'));
