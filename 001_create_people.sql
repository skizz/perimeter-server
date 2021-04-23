-- This is a sample migration.

create table sessions(
  id serial primary key,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

---- create above / drop below ----

drop table sessions;
