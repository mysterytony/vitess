#!/usr/bin/env bash

set -e

printf "\nStarting Vitess cluster\n"

export VTROOT=/vagrant
export VTDATAROOT=/tmp/vtdata-dev
export MYSQL_FLAVOR=MySQL56
cd "$VITESS_WORKSPACE"/examples/local
./zk-up.sh
./vtctld-up.sh --enable-grpc-static-auth

# bypage keyspace (2 shards)

SHARD='-80' KEYSPACE='bypage' UID_BASE='100' ./vttablet-up.sh --enable-grpc-static-auth
SHARD='80-' KEYSPACE='bypage' UID_BASE='200' ./vttablet-up.sh --enable-grpc-static-auth



# byuserid keyspace (4 shards)

SHARD='-40' KEYSPACE='byuserid' UID_BASE='300' ./vttablet-up.sh --enable-grpc-static-auth
SHARD='40-80' KEYSPACE='byuserid' UID_BASE='400' ./vttablet-up.sh --enable-grpc-static-auth
SHARD='80-c0' KEYSPACE='byuserid' UID_BASE='500' ./vttablet-up.sh --enable-grpc-static-auth
SHARD='c0-' KEYSPACE='byuserid' UID_BASE='600' ./vttablet-up.sh --enable-grpc-static-auth



./vtgate-up.sh --enable-grpc-static-auth

sleep 3

./lvtctl.sh InitShardMaster -force bypage/-80 test-100
./lvtctl.sh InitShardMaster -force bypage/80- test-200

./lvtctl.sh InitShardMaster -force byuserid/-40 test-300
./lvtctl.sh InitShardMaster -force byuserid/40-80 test-400
./lvtctl.sh InitShardMaster -force byuserid/80-c0 test-500
./lvtctl.sh InitShardMaster -force byuserid/c0- test-600

./lvtctl.sh ApplySchema -sql "$(cat create_bypage_table.sql)" bypage
./lvtctl.sh ApplyVSchema -vschema_file bypage_vschema.json bypage

./lvtctl.sh ApplySchema -sql "$(cat create_byuserid_table.sql)" byuserid
./lvtctl.sh ApplyVSchema -vschema_file byuserid_vschema.json byuserid

./lvtctl.sh RebuildVSchemaGraph

printf "\nadding data to database.\n\n"

sleep 3

cd $VITESS_WORKSPACE
ulimit -n 10000
export MYSQL_FLAVOR=MySQL56
export VT_MYSQL_ROOT=/usr
source dev.env
source /vagrant/dist/grpc/usr/local/bin/activate
./examples/local/client.sh

printf "\nadd schemaops user to db\n\n"

mysql --socket /tmp/vtdata-dev/vt_0000000100/mysql.sock -u root -e "create user 'schemaops'@'%' IDENTIFIED BY 'schemaops';"
mysql --socket /tmp/vtdata-dev/vt_0000000100/mysql.sock -u root -e "grant alter, create, delete, drop, index, insert, lock tables, select, trigger, update, super, replication client, replication slave  on *.* to 'schemaops'@'%';"

mysql --socket /tmp/vtdata-dev/vt_0000000200/mysql.sock -u root -e "create user 'schemaops'@'%' IDENTIFIED BY 'schemaops';"
mysql --socket /tmp/vtdata-dev/vt_0000000200/mysql.sock -u root -e "grant alter, create, delete, drop, index, insert, lock tables, select, trigger, update, super, replication client, replication slave  on *.* to 'schemaops'@'%';"

mysql --socket /tmp/vtdata-dev/vt_0000000300/mysql.sock -u root -e "create user 'schemaops'@'%' IDENTIFIED BY 'schemaops';"
mysql --socket /tmp/vtdata-dev/vt_0000000300/mysql.sock -u root -e "grant alter, create, delete, drop, index, insert, lock tables, select, trigger, update, super, replication client, replication slave  on *.* to 'schemaops'@'%';"

mysql --socket /tmp/vtdata-dev/vt_0000000400/mysql.sock -u root -e "create user 'schemaops'@'%' IDENTIFIED BY 'schemaops';"
mysql --socket /tmp/vtdata-dev/vt_0000000400/mysql.sock -u root -e "grant alter, create, delete, drop, index, insert, lock tables, select, trigger, update, super, replication client, replication slave  on *.* to 'schemaops'@'%';"

mysql --socket /tmp/vtdata-dev/vt_0000000500/mysql.sock -u root -e "create user 'schemaops'@'%' IDENTIFIED BY 'schemaops';"
mysql --socket /tmp/vtdata-dev/vt_0000000500/mysql.sock -u root -e "grant alter, create, delete, drop, index, insert, lock tables, select, trigger, update, super, replication client, replication slave  on *.* to 'schemaops'@'%';"

mysql --socket /tmp/vtdata-dev/vt_0000000600/mysql.sock -u root -e "create user 'schemaops'@'%' IDENTIFIED BY 'schemaops';"
mysql --socket /tmp/vtdata-dev/vt_0000000600/mysql.sock -u root -e "grant alter, create, delete, drop, index, insert, lock tables, select, trigger, update, super, replication client, replication slave  on *.* to 'schemaops'@'%';"
