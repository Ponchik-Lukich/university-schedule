FillDataQuery = """PRAGMA TablePathPrefix("{}");
DECLARE $roomsData AS List<Struct<
    id: Uint64,
    name: Utf8>>;

DECLARE $departmentsData AS List<Struct<
    id: Uint64,
    name: Utf8>>; 

DECLARE $groupsData AS List<Struct<
    id: Uint64,
    name: Utf8>>;

DECLARE $guestsData AS List<Struct<
    id: Uint64,
    short_name: Utf8>>; 

DECLARE $tutorsData AS List<Struct<
    id: Uint64,
    name: Utf8,
    short_name: Utf8>>; 

DECLARE $departmentLinksData AS List<Struct<
    id: Uint64,
    department_id: Uint64,
    semester: Uint32>>;

DECLARE $calendarPlanDepartmentLinksData AS List<Struct<
    calendar_plan_id: Uint64,
    department_link_id: Uint64>>;

DECLARE $calendarPlanGroupsData AS List<Struct<
    calendar_plan_id: Uint64,
    group_id: Uint64,
    choice: Uint64?>>;

REPLACE INTO rooms
SELECT
    id,
    name,
FROM AS_TABLE($roomsData);

REPLACE INTO departments
SELECT
    id,
    name,
FROM AS_TABLE($departmentsData);

REPLACE INTO groups
SELECT
    id,
    name,
FROM AS_TABLE($groupsData);

REPLACE INTO guests
SELECT
    id,
    short_name,
FROM AS_TABLE($guestsData);

REPLACE INTO tutors
SELECT
    id,
    name,
    short_name,
FROM AS_TABLE($tutorsData);

REPLACE INTO department_links
SELECT
    id,
    department_id,
    semester,
FROM AS_TABLE($departmentLinksData);

REPLACE INTO calendar_plan_department_links
SELECT
    calendar_plan_id,
    department_link_id,
FROM AS_TABLE($calendarPlanDepartmentLinksData);

REPLACE INTO calendar_plan_groups
SELECT
    calendar_plan_id,
    group_id,
    choice,
FROM AS_TABLE($calendarPlanGroupsData);
"""

FillDataQuery1 = """PRAGMA TablePathPrefix("{}");
DECLARE $CalendarPlanData AS List<Struct<
    id: Uint64,
    date: Utf8?,
    date_from: String?,
    date_to: String?,
    room_id: Uint32?,
    semester: Uint32,
    subject: Utf8,
    time_from: Uint32,
    time_to: Uint32,
    type: Utf8,
    week: Uint32?,
    week_day: Uint32?,>>;

REPLACE INTO calendar_plan
SELECT
    id,
    date,
    CAST(date_from AS Date) AS date_from,
    CAST(date_to AS Date) AS date_to,
    room_id,
    semester,
    subject,
    time_from,
    time_to,
    type,
    week,
    week_day,
FROM AS_TABLE($CalendarPlanData);
"""

FillDataQuery2 = """PRAGMA TablePathPrefix("{}");
DECLARE $guestsTimetableData AS List<Struct<
    id: Uint64,
    guest_id: Uint64>>; 

DECLARE $tutorsTimetableData AS List<Struct<
    id: Uint64,
    tutor_id: Uint64>>;

DECLARE $calendarPlanTutorsGuestsData AS List<Struct<
    id: Uint64,
    calendar_plan_id: Uint64,
    tutor_id: Uint64?,
    guest_id: Uint64?>>;

REPLACE INTO guests_timetable
SELECT
    id,
    guest_id,
FROM AS_TABLE($guestsTimetableData);

REPLACE INTO tutors_timetable
SELECT
    id,
    tutor_id,
FROM AS_TABLE($tutorsTimetableData);
"""