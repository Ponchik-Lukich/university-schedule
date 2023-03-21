# MEPhI Home Portal Schedule Tracker

## Project Overview

This project aims to create a serverless database using Yandex Database (YDB) to store the parsed data from the MEPhI Home Portal, which includes information about professors, groups, schedules, and more. A Yandex Cloud Function will be utilized to check for schedule changes every 5-10 minutes. If changes are detected, the information in the database will be updated, and a notification will be displayed to indicate recent updates. YDB uses its own SQL dialect called "YQL," but the project can also be replicated using PostgreSQL.

Initially, the plan was to create a website without user authentication, displaying all schedules along with recent changes if any. Additionally, a Telegram bot would be developed to track specific groups and send notifications about schedule changes or daily class notifications.

However, for the course requirements, user authentication can be added through the Home MEPhI portal, enabling fully authorized users (students) with full names and other relevant information.

## Features

- Serverless Yandex Database (YDB) to store parsed data from the MEPhI Home Portal
- Yandex Cloud Function to check for schedule changes every 5-10 minutes
- Notification system for recent updates
- (Optional) User authentication through Home MEPhI portal
- Telegram bot for tracking specific groups and sending notifications

## Installation

1. Clone the repository

`git clone https://github.com/yourusername/mephi-schedule-tracker.git`

2. Install dependencies

`cd mephi-schedule-tracker`
`pip install -r requirements.txt`

The remaining items are still in development:

3. Configure YDB and Yandex Cloud Function as per the [official documentation](https://cloud.yandex.com/docs)

4. Deploy the Yandex Cloud Function and set up the necessary triggers

5. Set up the Telegram bot by following the [official guide](https://core.telegram.org/bots)

6. Add the bot token and user authentication details to the configuration file

7. Run the application

## License

[MIT](https://choosealicense.com/licenses/mit/)

