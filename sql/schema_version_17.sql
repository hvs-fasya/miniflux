create table filters (
    id serial not null,
    user_id int not null,
    filter_name varchar,
    filters varchar ARRAY,
    created_at timestamp with time zone default now(),
    primary key (id),
    unique (filter_name),
    foreign key (user_id) references users(id) on delete cascade
);

create index filters_user_id_idx on filters using btree(user_id);
create index filters_filter_name_idx on filters using btree(filter_name);
