create table countries (
    id bigserial not null,
    name text not null,
    primary key (id),
    alpha3 text not null,
    unique (name)
);

create table headlines (
     id bigserial not null,
     hash text not null,
     published_at timestamp with time zone not null,
     title text not null,
     url text,
     content text,
     country_id bigint not null,
     visatype text,
     category_id bigint not null,
     icon_id bigint,
     primary key (id),
     unique (hash),
     foreign key (country_id) references countries(id) on delete cascade,
     foreign key (category_id) references categories(id) on delete cascade,
     foreign key (icon_id) references icons(id) on delete cascade
 );

create index headlines_content_idx on headlines using btree(content);
create index headlines_title_idx on headlines using btree(title);