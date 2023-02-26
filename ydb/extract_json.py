import json
import uuid


def save_json(filename, data):
    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=4)


def read_json(filename):
    with open(filename, 'r', encoding='utf-8') as f:
        return json.load(f)


def open_json(name):
    with open(name, encoding="utf8") as json_file:
        json_data = json.load(json_file)
        json_file.close()
        return json_data


def query_execute(query_head, query_body):
    s = query_head + '\n' + query_body + ';'
    print(s)


def fill_tutors():
    json_data = open_json('sources/postgres_public_tutors.json')
    count = 0
    head = """UPSERT INTO tutors (id, name, short_name, tutor_id) VALUES"""
    body = ""
    for tutor in json_data:
        body += '\t("' + str(uuid.uuid4()) + '", "' + tutor["name"] + '", "' + tutor["short_name"] + '", ' + str(tutor["id"]) + ')'
        if count != len(json_data) - 1 and count % 100 != 0:
            body += ',\n'
        else:
            query_execute(head, body)
            body = """"""
        count += 1


fill_tutors()
