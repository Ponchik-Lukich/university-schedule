import json
import ydb.iam
import os
import uuid
from dotenv import load_dotenv
from ydb.dbapi import OperationalError
# import classes from data.py
from data import *

load_dotenv()

def get_rooms_data():
    return [
        Entity(3, "room 1"),
        Entity(4, "room 2"),
    ]


FillDataQuery = """PRAGMA TablePathPrefix("{}");
DECLARE $roomsData AS List<Struct<
    id: Uint64,
    name: Utf8>>;
REPLACE INTO rooms
SELECT
    id,
    name,
FROM AS_TABLE($roomsData);
"""

driver_config = {
    "endpoint": os.getenv('ENDPOINT'),
    "database": os.getenv('DATABASE'),
    "credentials": ydb.iam.ServiceAccountCredentials.from_file(
        "./authorized_key.json"
    ),
    "root_certificates": ydb.load_ydb_root_certificate()
}


def executeScriptsFromFile(filename):
    # Open and read the file as a single buffer
    fd = open(filename, 'r', encoding="utf8")
    sqlFile = fd.read()
    fd.close()
    return sqlFile


def query_execute(query_head, query_body, session):
    s = query_head + '\n' + query_body + ';'
    # print(s)
    # print('\n')
    session.transaction().execute(
        s.format(ydb.iam.ServiceAccountCredentials.from_file("./authorized_key.json")),
        commit_tx=True,
    )


def readScriptsFromFile(filename):
    sql_data = executeScriptsFromFile(filename)
    sql_data = sql_data.split('\n')
    for sql in sql_data:
        try:
            print(sql)
        except OperationalError as msg:
            print("Command skipped: ", msg)


def parse_sql(head, filename, session):
    sql_data = executeScriptsFromFile(filename)
    sql_data = sql_data.split('\n')
    counter = 0
    body = ""
    for sql in sql_data:
        body += sql
        if counter != len(sql_data) - 1 and counter % 100 != 0 and counter != 0:
            body += '\n'
        else:
            query_execute(head, body[:-1], session)
            body = """"""
        counter += 1


with ydb.Driver(**driver_config) as driver:
    try:
        driver.wait(fail_fast=True, timeout=10)
        session = driver.table_client.session().create()
        # command = "insert into tutors (id, name, short_name) values"
        # parse_sql(command, 'sources/sql/tutors.sql', session)

        # command = "insert into tutors_timetable (id, tutor_id) values"
        # parse_sql(command, 'sources/sql/tutors_timetable.sql', session)
        prepared_query = session.prepare(
            FillDataQuery.format(driver_config['database']))
        session.transaction(ydb.SerializableReadWrite()).execute(
            prepared_query,
            {
                "$roomsData": get_rooms_data(),
            },
            commit_tx=True,
        )

        # command = "insert into tutors_lessons (id, tutor_id) values"
        # parse_sql(command, 'sources/sql/tutors_timetable.sql', session)
        exit(1)
    except TimeoutError:
        print("Connect failed to YDB")
        print("Last reported errors by discovery:")
        print(driver.discovery_debug_details())
        exit(1)
