create table users
(
    id       serial primary key,
    username varchar(255)
);

create table segments
(
    id         serial primary key,

    slug       varchar(255) not null,

    created_at timestamp default current_timestamp,
    deleted_at timestamp
);

create table user_segment
(
    id         serial primary key,

    user_id    int not null,
    segment_id int not null,

    end_at     timestamp,

    foreign key (segment_id) references segments (id),

    unique (user_id, segment_id, end_at)
);

create type operation_type as enum ('grant', 'revoke');

create table user_segment_history
(
    id         serial primary key,
    user_id    int            not null,
    segment_id int            not null,

    operation  operation_type not null,

    created_at timestamp      not null default current_timestamp,

    foreign key (user_id) references users (id),
    foreign key (segment_id) references segments (id)
);