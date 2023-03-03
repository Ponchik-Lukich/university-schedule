import json
from datetime import datetime


def get_json_data(file):
    with open(file, encoding="Utf8") as f:
        d = json.load(f)
        return d


class Entity(object):
    ___slots__ = ('id', 'name')

    def __init__(self, id, name):
        self.id = id
        self.name = name


class Guest(object):
    ___slots__ = ('id', 'short_name')

    def __init__(self, id, short_name):
        self.id = id
        self.short_name = short_name


class Tutor(object):
    ___slots__ = ('id', 'name', 'short_name')

    def __init__(self, id, name, short_name):
        self.id = id
        self.name = name
        self.short_name = short_name


class Department_link(object):
    ___slots__ = ('id', 'department_id', 'semester')

    def __init__(self, id, department_id, semester):
        self.id = id
        self.department_id = department_id
        self.semester = semester


class Calendar_plan(object):
    ___slots__ = (
        'id', 'date', 'date_from', 'date_to', 'room_id', 'semester', 'subject', 'time_from', 'time_to', 'type', 'week', 'week_day')

    def __init__(self, id, date, date_from, date_to, room_id, semester, subject, time_from, time_to, type, week, week_day):
        self.id = id
        self.room_id = room_id
        self.time_from = time_from
        self.time_to = time_to
        self.type = type
        self.week = week
        self.subject = subject
        self.week_day = week_day
        self.date = date
        if date_from is not None:
            self.date_from = bytes(date_from, "utf8")
        else:
            self.date_from = date_from
        if date_to is not None:
            self.date_to = bytes(date_to, "utf8")
        else:
            self.date_to = date_to
        self.semester = semester


class Calendar_plan_Department_links(object):
    ___slots__ = ('calendar_plan_id', 'department_link_id')

    def __init__(self, calendar_plan_id, department_link_id):
        self.calendar_plan_id = calendar_plan_id
        self.department_link_id = department_link_id


class Calendar_plan_Groups(object):
    ___slots__ = ('calendar_plan_id', 'group_id', 'choice')

    def __init__(self, calendar_plan_id, group_id, choice):
        self.calendar_plan_id = calendar_plan_id
        self.group_id = group_id
        self.choice = choice


class Guests_timetable(object):
    ___slots__ = ('id', 'guest_id')

    def __init__(self, id, guest_id):
        self.id = id
        self.guest_id = guest_id


class Tutors_timetable(object):
    ___slots__ = ('id', 'tutor_id')

    def __init__(self, id, tutor_id):
        self.id = id
        self.tutor_id = tutor_id


class Calendar_plan_Tutors_Guests(object):
    ___slots__ = ('id', 'calendar_plan_id', 'tutor_id', 'guest_id')

    def __init__(self, id, calendar_plan_id, tutor_id, guest_id):
        self.id = id
        self.calendar_plan_id = calendar_plan_id
        self.tutor_id = tutor_id
        self.guest_id = guest_id


def get_rooms_data():
    rooms_data = []
    rooms = get_json_data("./sources/json/rooms.json")
    for room in rooms:
        rooms_data.append(Entity(room["id"], room["name"]))
    return rooms_data


def get_departments_data():
    departments_data = []
    departments = get_json_data("./sources/json/departments.json")
    for department in departments:
        departments_data.append(Entity(department["id"], department["name"]))
    return departments_data


def get_groups_data():
    groups_data = []
    groups = get_json_data("./sources/json/departments.json")
    for group in groups:
        groups_data.append(Entity(group["id"], group["name"]))
    return groups_data


def get_guests_data():
    guests_data = []
    guests = get_json_data("./sources/json/guests.json")
    for guest in guests:
        guests_data.append(Guest(guest["id"], guest["short_name"]))
    return guests_data


def get_tutors_data():
    tutors_data = []
    tutors = get_json_data("./sources/json/tutors.json")
    for tutor in tutors:
        tutors_data.append(Tutor(tutor["id"], tutor["name"], tutor["short_name"]))
    return tutors_data


def get_department_links_data():
    department_links_data = []
    department_links = get_json_data("./sources/json/new_department_links.json")
    for department_link in department_links:
        department_links_data.append(Department_link(department_link["id"], department_link["department_id"], department_link["semester"]))
    print(department_links_data)
    return department_links_data


def get_calendar_plan_data():
    calendar_plan_data = []
    calendar_plans = get_json_data("./sources/json/calendar_plan.json")
    for calendar_plan in calendar_plans:
        time_from = calendar_plan["time_from"]
        time_to = calendar_plan["time_to"]
        int_time_from = int(time_from[:2]) * 60 + int(time_from[3:5])
        int_time_to = int(time_to[:2]) * 60 + int(time_to[3:5])
        calendar_plan_data.append(
            Calendar_plan(calendar_plan["id"], calendar_plan["date"],
                          calendar_plan["date_from"], calendar_plan["date_to"], calendar_plan["room_id"], calendar_plan["semester"],
                          calendar_plan["subject"], int_time_from,
                          int_time_to, calendar_plan["type"], calendar_plan["week"], calendar_plan["week_day"]))
    return calendar_plan_data


def get_calendar_plan_department_links_data():
    calendar_plan_department_links_data = []
    calendar_plan_department_links = get_json_data("./sources/json/calendar_plan_department_links.json")
    for calendar_plan_department_link in calendar_plan_department_links:
        calendar_plan_department_links_data.append(Calendar_plan_Department_links(calendar_plan_department_link
                                                                                  ["calendar_plan_id"],
                                                                                  calendar_plan_department_link
                                                                                  ["department_link_id"]))
    return calendar_plan_department_links_data


def get_calendar_plan_groups_data():
    calendar_plan_groups_data = []
    calendar_plan_groups = get_json_data("./sources/json/calendar_plan_groups.json")
    for calendar_plan_group in calendar_plan_groups:
        calendar_plan_groups_data.append(Calendar_plan_Groups(calendar_plan_group["calendar_plan_id"],
                                                              calendar_plan_group["group_id"],
                                                              calendar_plan_group["choice"]))
    return calendar_plan_groups_data


def get_guests_timetable_data():
    guests_timetable_data = []
    guests_timetable = get_json_data("./sources/json/guests_timetable.json")
    for guest_timetable in guests_timetable:
        guests_timetable_data.append(Guests_timetable(guest_timetable["id"], guest_timetable["guest_id"]))
    return guests_timetable_data


def get_tutors_timetable_data():
    tutors_timetable_data = []
    tutors_timetable = get_json_data("./sources/json/tutors_timetable.json")
    for tutor_timetable in tutors_timetable:
        tutors_timetable_data.append(Tutors_timetable(tutor_timetable["id"], tutor_timetable["tutor_id"]))
    return tutors_timetable_data


def get_calendar_plan_tutors_guests_data():
    calendar_plan_tutors_guests_data = []
    calendar_plan_tutors_guests = get_json_data("./sources/json/calendar_plan_tutors_guests.json")
    for calendar_plan_tutor_guest in calendar_plan_tutors_guests:
        calendar_plan_tutors_guests_data.append(Calendar_plan_Tutors_Guests(calendar_plan_tutor_guest["id"],
                                                                            calendar_plan_tutor_guest[
                                                                                "calendar_plan_id"],
                                                                            calendar_plan_tutor_guest["tutor_id"],
                                                                            calendar_plan_tutor_guest["guest_id"]))
    return calendar_plan_tutors_guests_data
