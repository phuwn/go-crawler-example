
-- +migrate Up

create table pancake_squad (
	id numeric not null unique primary key,
	image json,
	attributes json
);

-- +migrate Down
drop table pancake_squad;