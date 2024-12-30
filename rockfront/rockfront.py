import reflex as rx


import logging
import requests
import urllib.parse

from rockfront.sidebar import sidebar
from rockfront.navbar import navbar
# from rockbackup.policy import add_policy_page



backend = "http://localhost:8000"


class ValidationError(Exception):
    pass


class PolicyController:
    def __init__(self):
        self.request = requests
        self.base_url = urllib.parse.urljoin(backend, "service/")

    def add_file(self, path, host, retention, repository, full_bk_day, start_time):
        url = urllib.parse.urljoin(self.base_url, "file/open")
        full_cron, incr_cron = self.conv_param_to_cron(full_bk_day, start_time)

        payload = dict(
            source_path=path,
            hostname=host,
            retention=retention,
            repository_id=repository,
            full_backup_schedule=full_cron,
            incr_backup_schedule=incr_cron,
            start_time=start_time,
        )
        logging.info("open file service payload: %s", payload)
        resp = self.request.post(url, json=payload)
        return resp

    @staticmethod
    def conv_param_to_cron(full_day, start_time):
        full_cron, incr_cron = None, None

        incr_weekdays = [1, 2, 3, 4, 5, 6, 7]
        incr_weekdays.remove(full_day)

        hour, minute = start_time.split(":")

        full_cron = f"{minute} {hour} {full_day} * *"
        incr_cron = f"{minute} {hour} {','.join(list(map(str, incr_weekdays)))} * *"

        return full_cron, incr_cron


policy_controller = PolicyController()


class Requests:
    def get(self, url):
        resp = requests.get(url)
        return resp.json()

    def post(self, url, payload):
        resp = requests.post(url, payload)
        return resp.json()



def validate_policy(data):
    # breakpoint()
    musts = ["hostname", "source_path", "retention", "repository"]
    for m in musts:
        if m not in data and data[m] is None:
            raise ValidationError(f"{m} is empty")



def show_table(header, rows):

    header_rows = [rx.table.column_header_cell(x) for x in header]
    rx_header = rx.table.header(rx.table.row(*header_rows))

    body_rows = []

    for row in rows:
        table_rows = []
        for item in row:
            table_rows.append(rx.table.cell(item))
        body_rows.append(rx.table.row(*table_rows))

    return rx.table.root(
        rx_header,
        rx.table.body(*body_rows),
    )


def show_policy():
    parts = "service/file/get"
    url = urllib.parse.urljoin(backend, parts)
    # print("Geting response of url: %s" % url)
    ps = requests.get(url)

    header = [
        "Policy ID",
        "Source Type",
        "Source Path",
        "Source host",
        "Hostname",
        "Retention",
        "Status",
    ]
    table = []

    # if ps.status_code != 200:
    #     put_text("error")
    #     # return

    # put_text(ps.json())
    content = ps.json()
    if content is None:
        content = []

    for line in content:
        row = [
            line["id"],
            line["source_type"],
            line["source_path"],
            line["source_host"],
            line["hostname"],
            line["retention"],
            line["status"],
        ]
        table.append(row)

    return show_table(header, table)


def index() -> rx.Component:
    return rx.center(
        rx.vstack(
            navbar(),
            rx.hstack(
                sidebar(),
                rx.hstack(
                    # add_policy_page,
                    rx.vstack(
                        show_policy(),
                        align="center",
                    )
                )
            )
        )
    )


app = rx.App()
app.add_page(index)

# def show_file_policy_form():

#     data = input_group(
#         "Add Policy",
#         [
#             input("source path", name="source_path", other_html_attrs=dict(size=8, maxlength=10)),
#             select("host:", ["host1", "host2"], name="hostname"),
#             input("retention(days):", name="retention", type=NUMBER),
#             input("repository id:", name="repository", type=NUMBER),
#             select(
#                 "full backup day(week):", [1, 2, 3, 4, 5, 6, 7], name="full_backup_day"
#             ),
#             input("start time", name="start_time", type=TIME),
#             # input('Input your age', name='age', type=NUMBER, validate=check_age)
#         ],
#         validate=validate_policy,
#     )
#     result = policy_controller.add_file(
#         data["source_path"],
#         data["hostname"],
#         data["retention"],
#         data["repository"],
#         data["full_backup_day"],
#         data["start_time"],
#     )
#     breakpoint()
#     logging.info("response status code: %s", result.status_code)
#     logging.info("response: %s", result.json())

#     toast("Policy Added successfully")


# def show_policy():
#     parts = "service/file/get"
#     url = urllib.parse.urljoin(backend, parts)
#     # print("Geting response of url: %s" % url)
#     ps = requests.get(url)

#     header = [
#         "policy id",
#         "source_type",
#         "source path",
#         "source_host",
#         "hostname",
#         "retention",
#         "status",
#     ]
#     table = []

#     # breakpoint()
#     if ps.status_code != 200:
#         put_text("error")
#         # return

#     put_text(ps.json())
#     content = ps.json()
#     if content is None:
#         content = []

#     for line in content:
#         row = [
#             line["id"],
#             line["source_type"],
#             line["source_path"],
#             line["source_host"],
#             line["hostname"],
#             line["retention"],
#             line["status"],
#         ]
#         table.append(row)

#     put_table(table, header)
