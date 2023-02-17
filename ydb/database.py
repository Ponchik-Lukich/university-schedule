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

def parse_departments(session):
    with open("sources/department_timetable.json", encoding="utf8") as json_file:
        json_data = json.load(json_file)
        json_file.close()
        s = """UPSERT INTO departments (id, name) VALUES\n"""
        count = 0
        for semester in json_data:
            for department in json_data[semester]:
                s += '\t("' + str(uuid.uuid4()) + '", "'
                s += json_data[semester][department]['name']
                s += '")'
                s += ',\n'
        s += ';'
        print(s)
        # session.transaction().execute(
        #     s.format(ydb.iam.ServiceAccountCredentials.from_file("./authorized_key.json")),
        #     commit_tx=True,
        # )
        # session.transaction().execute(
        #     """
        #     UPSERT INTO departments (id, name) VALUES
        #         ("b44babeb-e72a-4396-9fd6-fed3572d2d9f", "Бебра");
        #     """.format(ydb.iam.ServiceAccountCredentials.from_file("./authorized_key.json")),
        #     commit_tx=True,
        # )
        # print(s)
        # session.transaction().execute(
        #     s.format(os.getenv('DATABASE')),
        #     commit_tx=True,
        # )
        # print(s)
        # print(json_data['2'])

def query_execute(query_head, query_body, session):
    s = query_head + '\n' + query_body + ';'
    session.transaction().execute(
        s.format(ydb.iam.ServiceAccountCredentials.from_file("./authorized_key.json")),
        commit_tx=True,
    )

def parse_departments_names(session):
    with open("sources/department_timetable.json", encoding="utf8") as json_file:
        json_data = json.load(json_file)
        json_file.close()
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

def parse_groups(session):
    with open("sources/groups.json", encoding="utf8") as json_file:
        json_data = json.load(json_file)
        json_file.close()
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



with ydb.Driver(**driver_config) as driver:
    try:
        driver.wait(fail_fast=True, timeout=1)
        session = driver.table_client.session().create()
        # parse_departments(session)
        # session.transaction().execute(
        #     """
        #     UPSERT INTO departments (id, name) VALUES
        #         ("2e010b71-1531-47b7-83a5-875883fda9a7", "Кафедра «Технологии замкнутого ядерного топливного цикла» (89)");
        #     """.format(ydb.iam.ServiceAccountCredentials.from_file("./authorized_key.json")),
        #     commit_tx=True,
        # )
        parse_departments_names(session)


        # session.create_table(driver_config['database'] + '/test',
        #                      ydb.TableDescription()
        #                      .with_column(ydb.Column('series_id', ydb.PrimitiveType.Uint64))  # not null column
        #                      .with_column(ydb.Column('title', ydb.OptionalType(ydb.PrimitiveType.Utf8)))
        #                      .with_column(ydb.Column('series_info', ydb.OptionalType(ydb.PrimitiveType.Utf8)))
        #                      .with_column(ydb.Column('release_date', ydb.OptionalType(ydb.PrimitiveType.Uint64)))
        #                      .with_primary_key('series_id'))
        exit(1)
    except TimeoutError:
        print("Connect failed to YDB")
        print("Last reported errors by discovery:")
        print(driver.discovery_debug_details())
        exit(1)

