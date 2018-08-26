create table alerts (
     id bigserial not null,
     last_updated timestamp with time zone not null,
     still_valid timestamp with time zone not null,
     latest_updates text,
     risk_level text,
     risk_details text,
     health_title text,
     health_date timestamp with time zone not null,
     health_content text,
     country_id bigint not null,
     primary key (id),
     foreign key (country_id) references countries(id) on delete cascade
 );
