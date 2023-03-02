import json
import ydb.iam
import os
import uuid
from dotenv import load_dotenv
from ydb.dbapi import OperationalError
from query import FillDataQuery, FillDataQuery1, FillDataQuery2
from data import *

load_dotenv()


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
                "$departmentsData": get_departments_data(),
                "$groupsData": get_groups_data(),
                "$guestsData": get_guests_data(),
                "$tutorsData": get_tutors_data(),
                "$departmentLinksData": get_department_links_data(),
                "$calendarPlanDepartmentLinksData": get_calendar_plan_department_links_data(),
                "$calendarPlanGroupsData": get_calendar_plan_groups_data(),
            },
            commit_tx=True,
        )
        prepared_query1 = session.prepare(
            FillDataQuery1.format(driver_config['database']))
        session.transaction(ydb.SerializableReadWrite()).execute(
            prepared_query1,
            {
                "$CalendarPlanData": get_calendar_plan_data(),
            },
            commit_tx=True,
        )
        prepared_query2 = session.prepare(
            FillDataQuery2.format(driver_config['database']))
        session.transaction(ydb.SerializableReadWrite()).execute(
            prepared_query2,
            {
                "$guestsTimetableData": get_guests_timetable_data(),
                "$tutorsTimetableData": get_tutors_timetable_data(),
                "$calendarPlanTutorsGuestsData": get_calendar_plan_tutors_guests_data(),
            },
            commit_tx=True,
        )

        exit(1)
    except TimeoutError:
        print("Connect failed to YDB")
        print("Last reported errors by discovery:")
        print(driver.discovery_debug_details())
        exit(1)
