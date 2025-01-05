package main

import (
	"fmt"
	"os"
	"os/signal"
	"rockbackup/backend/api"
	"rockbackup/backend/async"
	"rockbackup/backend/async/taskdef"
	"rockbackup/backend/db"
	"rockbackup/backend/keeperr"
	"rockbackup/backend/log"
	"rockbackup/backend/scheduler"
	"rockbackup/backend/schedules"
	"rockbackup/backend/service"
	"sync"

	gocron "github.com/robfig/cron/v3"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	AppHome     string
	LogBasePath string
	Config      *viper.Viper
	logger      *log.Logger
	cronLogger  *log.Logger
	RedisOpt    redis.Options
	WebEngine   *gin.Engine
	GormLog     string
	DB          *db.DB
	BackupSvc   service.BackupServiceI
)

func init() {
	initEnv()
	initLog()
	initConfig()
	initDB()
	initRedis()
	initCmd()
}

func initEnv() {
	AppHome = os.Getenv("rockbackup_home")
	if AppHome == "" {
		AppHome, _ = os.Getwd()
		fmt.Printf("env: rockbackup_home is not set, it is set default: %s\n", AppHome)
	}
}

func initLog() {
	// LogBasePath = filepath.Join(AppHome, "logs")
	// GormLog = filepath.Join(LogBasePath, "gorm.log")
	// cronLog := filepath.Join(LogBasePath, "cron.log")
	// backendPath := filepath.Join(LogBasePath, "backend.log")

	logger = log.New("backend.log")
	cronLogger = log.New("cron.log")
}

func initConfig() {
	Config = viper.New()
	Config.SetConfigName("config")
	Config.SetConfigType("toml")
	Config.AddConfigPath(".")

	err := Config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	fmt.Println("config is loaded")
}

func initDB() {
	DB = db.InitTest()
	// DB := db.InitDB(AppHome, Config, GormLog)
	if err := DB.AutoMigrate(); err != nil {
		panic(err)
	}
}

func initRedis() {
	addr := Config.GetString("redis.addr")
	db := Config.GetInt("redis.db")

	RedisOpt = redis.Options{
		Addr: addr,
		DB:   db,
	}
}

var cmdStartWorker = &cobra.Command{
	Use:   "worker",
	Short: "start asynq worker",
	Run: func(cmd *cobra.Command, args []string) {
		asynqRedisOpt := asynq.RedisClientOpt{Addr: RedisOpt.Addr, DB: RedisOpt.DB, Password: RedisOpt.Password}
		srv := asynq.NewServer(
			asynqRedisOpt,
			asynq.Config{
				Concurrency: 100,
				Queues: map[string]int{
					"critical": 6,
					"default":  3,
					"low":      1,
				},
			},
		)
		mux := asynq.NewServeMux()
		mux.HandleFunc(taskdef.TaskTypeBackupJobFile, async.MakeHandleBackupFileTask(Config, DB))

		if err := srv.Run(mux); err != nil {
			logger.Fatalf("could not run server: %v", err)
		}
	},
}

var cmdScheduler = &cobra.Command{
	Use:   "server",
	Short: "start serving",
	Run: func(cmd *cobra.Command, args []string) {

		keer := keeperr.NewKeepErr()
		// agentd := agentd.NewAgentd()
		// logbackup := logbackup.NewDBLogWatch(keer, agentd)
		asynqRedisOpt := asynq.RedisClientOpt{Addr: RedisOpt.Addr, DB: RedisOpt.DB, Password: RedisOpt.Password}
		asynqClient := asynq.NewClient(asynqRedisOpt)
		asyncHandler := scheduler.NewHandler(asynqClient)

		wg := &sync.WaitGroup{}
		JSched := scheduler.New(Config, DB, asyncHandler)

		gocronLogger := gocron.VerbosePrintfLogger(cronLogger)
		cron := gocron.New(gocron.WithLogger(gocronLogger))
		schedulesHandler := schedules.NewHandler(JSched)
		tSched := schedules.New(Config, DB, schedulesHandler, cron)

		BackupSvc := service.New(DB, tSched, JSched)
		webapi := api.New(BackupSvc)

		// wg.Add(1)
		// go func() {
		// 	defer wg.Done()
		// 	logbackup.Watch()
		// }()

		wg.Add(1)
		go func() {
			defer wg.Done()
			keer.Start()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			cron.Start()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			JSched.Start()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			tSched.Start()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			webapi.Start()
		}()

		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c
			// stop web api first
			webapi.Stop()
			tSched.Stop()
			JSched.Stop()
			cron.Stop()
			keer.Stop()
			fmt.Println("App is interrupted")
		}()

		wg.Wait()
	},
}

var cmdRoot = &cobra.Command{
	Use: "rock enterprise backup",
}

func initCmd() {
	cmdRoot.AddCommand(cmdStartWorker)
	cmdRoot.AddCommand(cmdScheduler)
}

func main() {
	if err := cmdRoot.Execute(); err != nil {
		panic(err)
	}
}
