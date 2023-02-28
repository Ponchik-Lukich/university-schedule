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
    ___slots__ = ('id', 'department_id')

    def __init__(self, id, department_id):
        self.id = id
        self.department_id = department_id


class Calendar_plan(object):
    ___slots__ = (
        'id', 'room_id', 'time_from', 'time_to', 'type', 'week', 'subject', 'week_day', 'date', 'date_from', 'date_to',
        'semester')

    def __init__(self, id, room_id, time_from, time_to, type, week, subject, week_day, date, date_from, date_to,
                 semester):
        self.id = id
        self.room_id = room_id
        self.time_from = time_from
        self.time_to = time_to
        self.type = type
        self.week = week
        self.subject = subject
        self.week_day = week_day
        self.date = date
        self.date_from = date_from
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


