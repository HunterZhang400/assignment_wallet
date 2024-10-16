package db_util

import (
	"assigement_wallet/src/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

// map[string]
var db *sqlx.DB

func InitDB() error {
	var err error
	db, err = sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s "+
		" host=%s port=%d dbname=%s sslmode=disable",
		config.ServerConfigs.Postgresql.UserName, config.ServerConfigs.Postgresql.Password,
		config.ServerConfigs.Postgresql.HostIP, config.ServerConfigs.Postgresql.Port,
		config.ServerConfigs.Postgresql.DbName))
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetDB() *sqlx.DB {
	return db
}

func InitDemoData(db *sqlx.DB) error {
	count := 0
	err := db.Get(&count, "SELECT count(1)  table_name FROM information_schema.tables"+
		" WHERE table_schema='public'"+
		"  AND table_type='BASE TABLE'")
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	db.MustExec(demoTableDataCreateSQL)
	return nil
}

const demoTableDataCreateSQL = `
CREATE TABLE wallet_balance (
    user_id text
    primary key,
    balance bigint
);
comment on table wallet_balance is 'user wallet balance';
comment on column wallet_balance.user_id is 'user unique and non-predictable id';
comment on column wallet_balance.balance is 'cents, like 525 represents $5.25 ';
alter table wallet_balance  owner to postgresql;

CREATE TABLE wallet_detail (
    id         bigint generated always as identity (cache 1000),
    user_id    text ,
    flow_type  smallint,
    amount     bigint,
    balance    bigint,
    occur_time timestamp not null,
    summary    text                   
);
comment on table wallet_detail is 'user wallet transaction details';
comment on column wallet_detail.id is 'transaction unique id';
comment on column wallet_detail.flow_type is 'money type by accounting, debit use 1, and credit use 2';
comment on column wallet_detail.amount is 'amount of money, like 525 represents $5.25 ';
comment on column wallet_detail.amount is 'the balance after transaction ';
comment on column wallet_detail.occur_time is 'transaction occur time';
comment on column wallet_detail.summary is 'transaction summary';
alter table wallet_detail  owner to postgresql;

insert into wallet_balance (user_id, balance) values ('a',9925);
insert into wallet_balance (user_id, balance) values ('b',10000);
insert into wallet_balance (user_id, balance) values ('c',75);

insert into wallet_detail (user_id, flow_type, amount, balance, occur_time, summary) values ('a',1, 10000,10000, '2024-10-13 17:00:00','deposit');
insert into wallet_detail (user_id, flow_type, amount, balance, occur_time, summary) values ('a',2, 75, 9925,'2024-10-13 17:00:01','transfer to c');
insert into wallet_detail (user_id, flow_type, amount, balance, occur_time, summary) values ('b',1, 5000, 5000,'2024-10-13 17:00:00','deposit');
insert into wallet_detail (user_id, flow_type, amount, balance, occur_time, summary) values ('b',1, 5000, 10000, '2024-10-13 17:00:00','deposit');
insert into wallet_detail (user_id, flow_type, amount, balance, occur_time, summary) values ('c',1, 75,75, '2024-10-13 17:00:01', 'received from a');

`
