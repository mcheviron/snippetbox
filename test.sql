use snippetbox;

create table
  users (
    id integer not null primary key auto_increment,
    name varchar(255) not null,
    email varchar(255) not null,
    hashed_password char(60) not null,
    created datetime not null
  );

alter table users add constraint users_uc_email unique (email);
