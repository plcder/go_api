create table users (
  id serial primary key,
  name text not null,
  role text not null,
  age integer not null
);