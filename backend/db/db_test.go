package db

// func TestInitTest(t *testing.T) {
// 	db := InitTest()
// 	assert.NotNil(t, db)
// }

// func TestAutoMigrate(t *testing.T) {
// 	db := InitTest()
// 	err := db.AutoMigrate()
// 	assert.Nil(t, err)
// }

// func TestGetPolicy(t *testing.T) {
// 	db := InitTest()
// 	db.AutoMigrate()

// 	// Create a policy to test retrieval
// 	p := policy.Policy{
// 		Name: "Test Policy",
// 		BackupSource: policy.BackupSource{
// 			Name: "Test Source",
// 		},
// 	}
// 	db.g.Create(&p)

// 	retrievedPolicy, err := db.GetPolicy(p.ID)
// 	assert.Nil(t, err)
// 	assert.Equal(t, p.Name, retrievedPolicy.Name)
// 	assert.Equal(t, p.BackupSource.Name, retrievedPolicy.BackupSource.Name)
// }

// func TestAddSchedulerJob(t *testing.T) {
// 	db := InitTest()
// 	db.AutoMigrate()

// 	job := &schedulerjob.Job{
// 		Name: "Test Job",
// 	}
// 	err := db.AddSchedulerJob(job)
// 	assert.Nil(t, err)

// 	var retrievedJob schedulerjob.Job
// 	db.g.First(&retrievedJob, job.ID)
// 	assert.Equal(t, job.Name, retrievedJob.Name)
// }
