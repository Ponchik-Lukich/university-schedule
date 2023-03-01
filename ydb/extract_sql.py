import json
import ydb.iam
import os
import uuid
from dotenv import load_dotenv
from ydb.dbapi import OperationalError
# import classes from data.py
from data import *

load_dotenv()


FillDataQuery = """PRAGMA TablePathPrefix("{}");
DECLARE $roomsData AS List<Struct<
    id: Uint64,
    name: Utf8>>;
   
DECLARE $departmentsData AS List<Struct<
    id: Uint64,
    short_name: Utf8>>; 
    
DECLARE $groupsData AS List<Struct<
    id: Uint64,
    short_name: Utf8>>;
    
DECLARE $guestsData AS List<Struct<
    id: Uint64,
    short_name: Utf8>>; 
    
DECLARE $tutorsData AS List<Struct<
    id: Uint64,
    name: Utf8,
    short_name: Utf8>>; 
    
DECLARE $departmentLinksData AS List<Struct<
    id: Uint64,
    department_id: Uint64>>;
    
DECLARE $CalendarPlanData AS List<Struct<
    id: Uint64,
    room_id: Uint64,
    time_from: Timestamp,
    time_to: Timestamp,
    type: Utf8,
    week: Uint32,
    subject: Utf8,
    week_day: Uint32,
    date: String,
    date_from: Date,
    date_to: Date,
    semester: Uint32>>;
    
DECLARE $calendarPlanDepartmentLinksData AS List<Struct<
    calendar_plan_id: Uint64,
    department_link_id: Uint64>>;
    
DECLARE $calendarPlanGroupsData AS List<Struct<
    calendar_plan_id: Uint64,
    group_id: Uint64
    choice: Uint64>>;
    
DECLARE $guestsTimetableData AS List<Struct<
    id: Uint64,
    guest_id: Uint64>>; 
    
DECLARE $tutorsTimetableData AS List<Struct<
    id: Uint64,
    tutor_id: Uint64>>;
    
DECLARE $calendarPlanTutorsGuestsData AS List<Struct<
    id: Uint64,
    calendar_plan_id: Uint64,
    tutor_id: Uint64,
    guest_id: Uint64>>;
    
REPLACE INTO rooms
SELECT
    id,
    name,
FROM AS_TABLE($roomsData);

REPLACE INTO guests
SELECT
    id,
    short_name,
FROM AS_TABLE($guestsData);
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
        prepared_query = session.prepare(
            FillDataQuery.format(driver_config['database']))
        session.transaction(ydb.SerializableReadWrite()).execute(
            prepared_query,
            {
                "$roomsData": get_rooms_data(),
                "$guestsData": get_guests_data(),
            },
            commit_tx=True,
        )

        exit(1)
    except TimeoutError:
        print("Connect failed to YDB")
        print("Last reported errors by discovery:")
        print(driver.discovery_debug_details())
        exit(1)
