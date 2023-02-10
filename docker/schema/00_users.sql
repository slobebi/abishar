create table if not exists users
(
  id bigserial primary key,
  username varchar(50) not null,
  password varchar(255) not null,
  is_admin boolean default false not null,
  salt varchar(255),
  created_time timestamp with time zone default now() not null,
  updated_time timestamp with time zone default now() not null
);

insert into users (username, password, is_admin, salt)
values ('adminOrdent', 'a47d9692589d1355e6356b99b08583806633d7ea', true, '4a5e388a-e758-4b99-8c3f-4a811f5053e0');
