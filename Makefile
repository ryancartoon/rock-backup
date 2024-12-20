swag:
	swag init -g api.go --dir ./backend/api/ -o ./backend/api/docs
	swag fmt -g ./backend/api/api.go
