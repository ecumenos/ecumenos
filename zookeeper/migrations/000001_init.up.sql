begin;

create table public.admins
(
  id            bigint primary key,
  created_at    timestamp(0) with time zone default current_timestamp not null,
  updated_at    timestamp(0) with time zone default current_timestamp not null,
  deleted_at    timestamp(0) with time zone,
  tombstoned    boolean not null default false,
  email         text not null,
  password_hash text not null
);
create unique index admins_email_uindex on admins (email);

create table public.admin_roles
(
  id               bigint primary key,
  created_at       timestamp(0) with time zone default current_timestamp not null,
  updated_at       timestamp(0) with time zone default current_timestamp not null,
  deleted_at       timestamp(0) with time zone,
  tombstoned       boolean not null default false,
  name             text not null,
  permissions      jsonb not null,
  creator_admin_id bigint references admins (id) not null
);
create unique index admin_roles_name_uindex on admin_roles (name);

create table public.admins_admin_roles_relations
(
  receiver_admin_id bigint references admins (id) not null,
  granter_admin_id  bigint references admins (id),
  role_id           bigint references admin_roles (id) not null,
  granted_at        timestamp(0) with time zone default current_timestamp not null
);

create table public.admin_sessions
(
  id               bigint primary key,
  created_at       timestamp(0) with time zone default current_timestamp not null,
  updated_at       timestamp(0) with time zone default current_timestamp not null,
  expired_at       timestamp(0) with time zone not null,
  deleted_at       timestamp(0) with time zone,
  tombstoned       boolean not null default false,
  admin_id         bigint references admins (id) not null,
  token            text not null,
  refresh_token    text not null
);

create table public.orbes_socii
(
  id                bigint primary key,
  created_at        timestamp(0) with time zone default current_timestamp not null,
  updated_at        timestamp(0) with time zone default current_timestamp not null,
  deleted_at        timestamp(0) with time zone,
  tombstoned        boolean not null default false,
  owner_email       text not null,
  alive             boolean not null default false,
  robustness_status integer check (robustness_status >= 0) not null,
  last_pinged_at    timestamp(0) with time zone,
  region            text not null,
  name              text not null,
  description       text not null,
  url               text not null,
  api_key           text not null
);

create table public.orbes_socii_stats
(
  id              bigint primary key,
  created_at      timestamp(0) with time zone default current_timestamp not null,
  orbis_socius_id bigint references orbes_socii(id),
  alive           boolean not null default false
);

create table public.orbes_socii_launch_requests
(
  id                       bigint primary key,
  created_at               timestamp(0) with time zone default current_timestamp not null,
  email                    text not null,
  region                   text not null,
  orbis_socius_name        text not null,
  orbis_socius_description text not null,
  orbis_socius_url         text not null,
  status                   integer check (status >= 0) not null
);

create table public.orbes_socii_launch_invites
(
  id                              bigint primary key,
  created_at                      timestamp(0) with time zone default current_timestamp not null,
  email                           text not null,
  code                            text not null,
  api_key                         text not null,
  used                            boolean not null default false,
  orbis_socius_launch_requests_id bigint references orbes_socii_launch_requests(id),
  expired_at                      timestamp(0) with time zone not null
);

commit;
