
##### repo, backend, host

insert into backends ("name", "type", "path") values ("file1", "file-restic", "/tmp/rock-repo")

##### start service

go run main server
go run main worker
go run cmd/agent/main.go start


##### create policy

POST http://localhost:8000/service/file/open HTTP/1.1
content-type: application/json

{
  "source_path": "/home/ryan/codes/backup-source-path",
  "source_name": "rock-file-backup-1",
  "hostname": "localhost",
  "backup_plan": 1,
  "retention": 1,
  "full_backup_schedule": "0 0 * * *",
  "incr_backup_schedule": "0 12 * * *",
  "log_backup_schedule": "1",
  "start_time": "14:00",
  "backend_id": 1,
  "duration": 24
}


##### start a backup job

go run cmd/admin/main.go job add -t backup -b full -p 1