




```json
POST http://localhost:8000/service/file/open HTTP/1.1
content-type: application/json

{
    "source_path": "/home/ryan/codes/rock-backup/logs",
    "hostname": "localhost",
    "backup_plan":1,
    "full_backup_schedule": "*/10 * * * *",
    "incr_backup_schedule": "*/5 * * * *",
    // "log_backup_schedule": "*/2 * * * *",
    "retention" : 31,
    "backup_cycle": 7,
    "start_time": "20:00",
    "duration": 60,
    "repository_id": 1,
    "backup_cycle": 7
}



GET http://localhost:8000/service/file/get HTTP/1.1

{}

```
