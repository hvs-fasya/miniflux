create table alerts (
     id bigserial not null,
     last_updated timestamp with time zone not null,
     still_valid timestamp with time zone not null,
     latest_updates text,
     risk_level text,
     risk_details text,
     country_id bigint not null,
     primary key (id),
     foreign key (country_id) references countries(id) on delete cascade
 );

 create table healths (
     id bigserial not null,
     health_title text not null,
     health_link text,
     health_content text,
     last_updated timestamp with time zone not null,
     primary key (id),
     unique (health_title)
 );

 create table alert_health (
    country_id bigint not null,
    health_id bigint not null,
    alert_health_date timestamp with time zone not null,
    foreign key (country_id) references countries(id) on delete cascade,
    foreign key (health_id) references healths(id) on delete cascade
 );
