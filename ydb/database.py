# from yandex.cloud import ydb
import json

import ydb.iam
import os
import uuid
from dotenv import load_dotenv

load_dotenv()

driver_config = {
    "endpoint": os.getenv('ENDPOINT'),
    "database": os.getenv('DATABASE'),
    # "credentials": ydb.AnonymousCredentials()
    "credentials": ydb.iam.ServiceAccountCredentials.from_file(
        "./authorized_key.json"
    ),
    "root_certificates": ydb.load_ydb_root_certificate()
}


def save_json(filename, data):
    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=4)


def read_json(filename):
    with open(filename, 'r', encoding='utf-8') as f:
        return json.load(f)


def query_execute(query_head, query_body, session):
    s = query_head + '\n' + query_body + ';'
    session.transaction().execute(
        s.format(ydb.iam.ServiceAccountCredentials.from_file("./authorized_key.json")),
        commit_tx=True,
    )


def open_json(name):
    with open(name, encoding="utf8") as json_file:
        json_data = json.load(json_file)
        json_file.close()
        return json_data


def parse_groups(session):
    json_data = open_json('sources/groups.json')
    count = 0
    head = """UPSERT INTO groups (id, name) VALUES"""
    body = ""
    for group in json_data:
        body += '\t("' + str(uuid.uuid4()) + '", "' + json_data[group] + '")'
        if count != len(json_data) - 1 and count % 100 != 0:
            body += ',\n'
        else:
            query_execute(head, body, session)
            body = """"""
        count += 1


def parse_departments_names(session):
    json_data = open_json('sources/department_timetable.json')
    names = []
    for semester in json_data:
        for department in json_data[semester]:
            names.append(json_data[semester][department]['name'].replace('"', "|"))
    names = list(set(names))

    count = 0
    head = """UPSERT INTO departments (id, name) VALUES"""
    body = ""

    for name in names:
        body += '\t("' + str(uuid.uuid4()) + '", "' + name + '")'
        if count != len(names) - 1 and count != len(names) / 2:
            body += ',\n'
        else:
            query_execute(head, body, session)
            body = """"""
        count += 1

def parse_department_links(session):
    result_sets = session.transaction(ydb.SerializableReadWrite()).execute(
        """
        SELECT 
        id, 
        name
        FROM departments;
        """.format(ydb.iam.ServiceAccountCredentials.from_file("./authorized_key.json")),
        commit_tx=True,
    )
    for row in result_sets[0].rows:
        print("id: ", row.id, ", name: ", row.name)
    # print(result_sets)



with ydb.Driver(**driver_config) as driver:
    try:
        driver.wait(fail_fast=True, timeout=1)
        session = driver.table_client.session().create()
        # parse_departments_names(session)
        # parse_groups(session)
        parse_department_links(session)

        # session.create_table(driver_config['database'] + '/test',
        #                      ydb.TableDescription()
        #                      .with_column(ydb.Column('series_id', ydb.PrimitiveType.Uint64))  # not null column
        #                      .with_column(ydb.Column('title', ydb.OptionalType(ydb.PrimitiveType.Utf8)))
        #                      .with_column(ydb.Column('series_info', ydb.OptionalType(ydb.PrimitiveType.Utf8)))
        #                      .with_column(ydb.Column('release_date', ydb.OptionalType(ydb.PrimitiveType.Uint64)))
        #                      .with_primary_key('series_id'))

        # parse_departments(session)
        # session.transaction().execute(
        #     """
        #     UPSERT INTO departments (id, name) VALUES
        #         ("2e010b71-1531-47b7-83a5-875883fda9a7", "Кафедра «Технологии замкнутого ядерного топливного цикла» (89)");
        #     """.format(ydb.iam.ServiceAccountCredentials.from_file("./authorized_key.json")),
        #     commit_tx=True,
        # )
        exit(1)
    except TimeoutError:
        print("Connect failed to YDB")
        print("Last reported errors by discovery:")
        print(driver.discovery_debug_details())
        exit(1)
