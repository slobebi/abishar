create table if not exists votes (
  id bigserial primary key,
  user_id bigint not null,
  movie_id bigint not null,
  created_time timestamp with time zone default now() not null,
  updated_time timestamp with time zone default now() not null,
  constraint votes_user_id_fk foreign key (user_id)
    references users(id),
  constraint votes_movie_id_fk foreign key (movie_id)
    references movies(id)
);
