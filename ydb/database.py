# from yandex.cloud import ydb

import ydb.iam
# from ydb.aio import driver


driver_config = {
    "endpoint": "grpcs://ydb.serverless.yandexcloud.net:2135",
    "database": "/ru-central1/b1gsopf4nf64thcp1nbp/etnnv6lqhm7e2vaqgk89",
    # "credentials": ydb.AnonymousCredentials()
    "credentials": ydb.iam.ServiceAccountCredentials.from_file(
        "./authorized_key.json"
    ),
    "root_certificates": ydb.load_ydb_root_certificate()
}
with ydb.Driver(**driver_config) as driver:
    try:
        driver.wait(fail_fast=True, timeout=1)
        session = driver.table_client.session().create()

        session.create_table(driver_config['database'] + '/test',
                             ydb.TableDescription()
                             .with_column(ydb.Column('series_id', ydb.PrimitiveType.Uint64))  # not null column
                             .with_column(ydb.Column('title', ydb.OptionalType(ydb.PrimitiveType.Utf8)))
                             .with_column(ydb.Column('series_info', ydb.OptionalType(ydb.PrimitiveType.Utf8)))
                             .with_column(ydb.Column('release_date', ydb.OptionalType(ydb.PrimitiveType.Uint64)))
                             .with_primary_key('series_id'))
    except TimeoutError:
        print("Connect failed to YDB")
        print("Last reported errors by discovery:")
        print(driver.discovery_debug_details())
        exit(1)

