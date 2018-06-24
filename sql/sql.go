// Code generated by go generate; DO NOT EDIT.
// 2018-06-24 14:38:57.146456533 +0300 MSK m=+0.001603959

package sql

var SqlMap = map[string]string{
	"schema_version_1": `create table schema_version (
    version text not null
);

create table users (
    id serial not null,
    username text not null unique,
    password text,
    is_admin bool default 'f',
    language text default 'en_US',
    timezone text default 'UTC',
    theme text default 'default',
    last_login_at timestamp with time zone,
    primary key (id)
);

create table sessions (
    id serial not null,
    user_id int not null,
    token text not null unique,
    created_at timestamp with time zone default now(),
    user_agent text,
    ip text,
    primary key (id),
    unique (user_id, token),
    foreign key (user_id) references users(id) on delete cascade
);

create table categories (
    id serial not null,
    user_id int not null,
    title text not null,
    primary key (id),
    unique (user_id, title),
    foreign key (user_id) references users(id) on delete cascade
);

create table feeds (
    id bigserial not null,
    user_id int not null,
    category_id int not null,
    title text not null,
    feed_url text not null,
    site_url text not null,
    checked_at timestamp with time zone default now(),
    etag_header text default '',
    last_modified_header text default '',
    parsing_error_msg text default '',
    parsing_error_count int default 0,
    primary key (id),
    unique (user_id, feed_url),
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (category_id) references categories(id) on delete cascade
);

create type entry_status as enum('unread', 'read', 'removed');

create table entries (
    id bigserial not null,
    user_id int not null,
    feed_id bigint not null,
    hash text not null,
    published_at timestamp with time zone not null,
    title text not null,
    url text not null,
    author text,
    content text,
    status entry_status default 'unread',
    primary key (id),
    unique (feed_id, hash),
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (feed_id) references feeds(id) on delete cascade
);

create index entries_feed_idx on entries using btree(feed_id);

create table enclosures (
    id bigserial not null,
    user_id int not null,
    entry_id bigint not null,
    url text not null,
    size int default 0,
    mime_type text default '',
    primary key (id),
    foreign key (user_id) references users(id) on delete cascade,
    foreign key (entry_id) references entries(id) on delete cascade
);

create table icons (
    id bigserial not null,
    hash text not null unique,
    mime_type text not null,
    content bytea not null,
    primary key (id)
);

create table feed_icons (
    feed_id bigint not null,
    icon_id bigint not null,
    primary key(feed_id, icon_id),
    foreign key (feed_id) references feeds(id) on delete cascade,
    foreign key (icon_id) references icons(id) on delete cascade
);
`,
	"schema_version_10": `drop table tokens;

create table sessions (
    id text not null,
    data jsonb not null,
    created_at timestamp with time zone not null default now(),
    primary key(id)
);`,
	"schema_version_11": `alter table integrations add column wallabag_enabled bool default 'f';
alter table integrations add column wallabag_url text default '';
alter table integrations add column wallabag_client_id text default '';
alter table integrations add column wallabag_client_secret text default '';
alter table integrations add column wallabag_username text default '';
alter table integrations add column wallabag_password text default '';`,
	"schema_version_12": `alter table entries add column starred bool default 'f';`,
	"schema_version_13": `create index entries_user_status_idx on entries(user_id, status);
create index feeds_user_category_idx on feeds(user_id, category_id);
`,
	"schema_version_14": `alter table integrations add column nunux_keeper_enabled bool default 'f';
alter table integrations add column nunux_keeper_url text default '';
alter table integrations add column nunux_keeper_api_key text default '';`,
	"schema_version_15": `alter table enclosures alter column size set data type bigint;`,
	"schema_version_16": `alter table entries add column comments_url text default '';`,
	"schema_version_17": `create table filters (
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
`,
	"schema_version_2": `create extension if not exists hstore;
alter table users add column extra hstore;
create index users_extra_idx on users using gin(extra);
`,
	"schema_version_3": `create table tokens (
    id text not null,
    value text not null,
    created_at timestamp with time zone not null default now(),
    primary key(id, value)
);`,
	"schema_version_4": `create type entry_sorting_direction as enum('asc', 'desc');
alter table users add column entry_direction entry_sorting_direction default 'asc';
`,
	"schema_version_5": `create table integrations (
    user_id int not null,
    pinboard_enabled bool default 'f',
    pinboard_token text default '',
    pinboard_tags text default 'miniflux',
    pinboard_mark_as_unread bool default 'f',
    instapaper_enabled bool default 'f',
    instapaper_username text default '',
    instapaper_password text default '',
    fever_enabled bool default 'f',
    fever_username text default '',
    fever_password text default '',
    fever_token text default '',
    primary key(user_id)
)
`,
	"schema_version_6": `alter table feeds add column scraper_rules text default '';
`,
	"schema_version_7": `alter table feeds add column rewrite_rules text default '';
`,
	"schema_version_8": `alter table feeds add column crawler boolean default 'f';
`,
	"schema_version_9": `alter table sessions rename to user_sessions;`,
}

var SqlMapChecksums = map[string]string{
	"schema_version_1":  "00b2fa9e945565625c93ef9d4242a8b6583dc3cd7edf38d2fc95c0f3f7b926ae",
	"schema_version_10": "8faf15ddeff7c8cc305e66218face11ed92b97df2bdc2d0d7944d61441656795",
	"schema_version_11": "dc5bbc302e01e425b49c48ddcd8e29e3ab2bb8e73a6cd1858a6ba9fbec0b5243",
	"schema_version_12": "a95abab6cdf64811fc744abd37457e2928939d999c5ef00d2bdd9398e16f32fb",
	"schema_version_13": "9073fae1e796936f4a43a8120ebdb4218442fe7d346ace6387556a357c2d7edf",
	"schema_version_14": "4622e42c4a5a88b6fe1e61f3d367b295968f7260ab5b96481760775ba9f9e1fe",
	"schema_version_15": "13ff91462bdf4cda5a94a4c7a09f757761b0f2c32b4be713ba4786a4837750e4",
	"schema_version_16": "9d006faca62fd7ab787f64aef0e0a5933d142466ec4cab0e096bb920d2797e34",
	"schema_version_17": "ac5603f064f48a1add8b0f5fbf8e48817021e6783d89738ba35b32eb4ac72f01",
	"schema_version_2":  "e8e9ff32478df04fcddad10a34cba2e8bb1e67e7977b5bd6cdc4c31ec94282b4",
	"schema_version_3":  "a54745dbc1c51c000f74d4e5068f1e2f43e83309f023415b1749a47d5c1e0f12",
	"schema_version_4":  "216ea3a7d3e1704e40c797b5dc47456517c27dbb6ca98bf88812f4f63d74b5d9",
	"schema_version_5":  "46397e2f5f2c82116786127e9f6a403e975b14d2ca7b652a48cd1ba843e6a27c",
	"schema_version_6":  "9d05b4fb223f0e60efc716add5048b0ca9c37511cf2041721e20505d6d798ce4",
	"schema_version_7":  "33f298c9aa30d6de3ca28e1270df51c2884d7596f1283a75716e2aeb634cd05c",
	"schema_version_8":  "9922073fc4032d8922617ec6a6a07ae8d4817846c138760fb96cb5608ab83bfc",
	"schema_version_9":  "de5ba954752fe808a993feef5bf0c6f808e0a4ced5379de8bec8342678150892",
}
