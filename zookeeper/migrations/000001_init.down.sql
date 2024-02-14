begin;

drop table if exists admin_sessions cascade;
drop table if exists admins_admin_roles_relations cascade;
drop table if exists admin_roles cascade;
drop table if exists admins cascade;

drop table if exists comptus_sessions cascade;
drop table if exists compti cascade;

drop table if exists orbes_socii_launch_invites cascade;
drop table if exists orbes_socii_launch_requests cascade;
drop table if exists orbes_socii_stats cascade;
drop table if exists orbes_socii cascade;

commit;
