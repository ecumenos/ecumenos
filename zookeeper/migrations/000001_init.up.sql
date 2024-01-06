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

commit;
