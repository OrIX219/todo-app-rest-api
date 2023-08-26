-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE users
(
    id            serial primary key,
    name          varchar(255) not null,
    username      varchar(255) not null,
    password_hash varchar(255) not null 
);

CREATE TABLE todo_lists
(
    id          serial primary key,
    title       varchar(255) not null,
    description varchar(255)
);


CREATE TABLE users_lists
(
    id    serial primary key,
    user_id int,
    list_id int,
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (list_id) references todo_lists(id) on delete cascade
);

CREATE TABLE todo_items
(
    id          serial primary key,
    title       varchar(255) not null,
    description varchar(255),
    done        boolean not null default false
);

CREATE TABLE lists_items
(
    id    serial primary key,
    item_id int,
    list_id int,
    foreign key (item_id) references todo_items(id) on delete cascade,
    foreign key (list_id) references todo_lists(id) on delete cascade
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE lists_items;

DROP TABLE users_lists;

DROP TABLE todo_lists;

DROP TABLE users;

DROP TABLE todo_items;

-- +goose StatementEnd
