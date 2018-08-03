create table visaupdates (
     id bigserial not null,
     hash text not null,
     published_at timestamp with time zone not null,
     title text not null,
     authority text,
     content text,
     country_id bigint not null,
     visatype text,
     primary key (id),
     unique (hash),
     foreign key (country_id) references countries(id) on delete cascade
 );

create index headlines_content_idx on headlines using btree(content);
create index headlines_title_idx on headlines using btree(title);
