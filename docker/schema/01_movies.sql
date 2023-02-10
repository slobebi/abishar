create table if not exists movies (
  id bigserial primary key,
  title varchar(255) not null,
  description text not null,
  duration integer not null,
  artists text not null,
  genres text not null,
  watch_url text not null,
  views bigint default 0 not null,
  created_time timestamp with time zone default now() not null,
  updated_time timestamp with time zone default now() not null
);
