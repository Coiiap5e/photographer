package repository

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Coiiap5e/photographer/model"
	"github.com/Coiiap5e/photographer/testutils"
	"github.com/stretchr/testify/suite"
)

func TestShootRepositoryTestSuit(t *testing.T) {
	suite.Run(t, new(ShootRepositoryTestSuit))
}

type ShootRepositoryTestSuit struct {
	suite.Suite
	ctx        context.Context
	db         *testutils.TestDB
	repo       *ShootRepository
	testClient *model.Client
}

func (suite *ShootRepositoryTestSuit) SetupSuite() {
	suite.ctx = context.Background()

	var err error

	suite.db, err = testutils.CreateTestDB(suite.ctx)
	suite.Require().NoError(err, "Failed to setup test database")
	suite.Require().NotNil(suite.db, "TestDB should not be nil")
	suite.Require().NotNil(suite.db.GetDB(), "DB connection should not be nil")

	suite.repo = NewShootRepository(suite.db.GetDB())
}

func (suite *ShootRepositoryTestSuit) SetupTest() {
	err := suite.db.CleanTables(suite.ctx)
	suite.Require().NoError(err)

	clientRepo := NewClientRepository(suite.db.GetDB())
	testClient := testutils.CreateTestClient()
	err = clientRepo.AddClient(suite.ctx, testClient)
	suite.Require().NoError(err)
	suite.testClient = testClient
}

func (suite *ShootRepositoryTestSuit) TearDownSuite() {
	if suite.db != nil {
		if err := suite.db.Cleanup(suite.ctx); err != nil {
			log.Fatalf("Failed to cleanup test database: %v", err)
		}
	}
}

func (suite *ShootRepositoryTestSuit) TestAddShoot() {
	// Given
	testShoot := testutils.CreateTestShoot(suite.testClient.Id)

	// When
	err := suite.repo.AddShoot(suite.ctx, testShoot)

	// Then
	suite.NoError(err, "should not return error")
	suite.Greater(testShoot.Id, 0, "should assign positive ID to shoot object")
	suite.NotZero(testShoot.CreatedAt, "should set creation timestamp in shoot object")
}

func (suite *ShootRepositoryTestSuit) TestGetShootByID() {
	suite.T().Run("should successfully get a shoot by ID", func(t *testing.T) {
		// Given
		testShoot := testutils.CreateTestShoot(suite.testClient.Id)
		err := suite.repo.AddShoot(suite.ctx, testShoot)
		suite.NoError(err, "shoot should be created")

		// When
		retrievedShoot, err := suite.repo.GetShootByID(suite.ctx, testShoot.Id)

		// Then
		suite.NoError(err, "should retrieve shoot without error")
		suite.NotNil(retrievedShoot, "should return shoot object")
		suite.Equal(testShoot.Id, retrievedShoot.Id, "should retrieve correct shoot by ID")
		suite.Equal(testShoot.ShootDate.Format("2006-01-02"), retrievedShoot.ShootDate.Format("2006-01-02"),
			"should retrieve correct shoot date")
		suite.Equal(testShoot.StartTime.Format("15:04"), retrievedShoot.StartTime.Format("15:04"),
			"should retrieve correct shoot start time")
		suite.Equal(testShoot.EndTime.Format("15:04"), retrievedShoot.EndTime.Format("15:04"),
			"should retrieve correct shoot end time")
		suite.Equal(testShoot.ShootPrice, retrievedShoot.ShootPrice, "should retrieve correct shoot price")
		suite.Equal(testShoot.ShootLocation, retrievedShoot.ShootLocation, "should retrieve correct shoot location")
		suite.Equal(testShoot.ShootType, retrievedShoot.ShootType, "should retrieve correct shoot type")
		suite.Equal(testShoot.Notes, retrievedShoot.Notes, "should retrieve correct notes")
	})
	suite.T().Run("should return error when get non-existent shoot", func(t *testing.T) {
		// Given
		nonExistentID := 9999

		// When
		retrievedShoot, err := suite.repo.GetShootByID(suite.ctx, nonExistentID)

		// Then
		suite.Error(err, "should return error for non-existent shoot")
		suite.Nil(retrievedShoot, "should not return shoot object")
	})

}

func (suite *ShootRepositoryTestSuit) TestDeleteShoot() {
	suite.T().Run("should successfully delete shoot", func(t *testing.T) {
		// Given
		testShoot := testutils.CreateTestShoot(suite.testClient.Id)
		err := suite.repo.AddShoot(suite.ctx, testShoot)
		suite.NoError(err, "shoot should be created")

		shootID := testShoot.Id

		// When
		err = suite.repo.DeleteShoot(suite.ctx, testShoot.Id)

		// Then
		suite.NoError(err, "should not return error")

		deletedShoot, err := suite.repo.GetShootByID(suite.ctx, shootID)
		suite.Error(err, "should return error when getting deleted shoot")
		suite.Nil(deletedShoot, "should not return deleted shoot")
	})
	suite.T().Run("should return error when deleting non-existent shoot", func(t *testing.T) {
		// Given
		nonExistentID := 9999

		// When
		err := suite.repo.DeleteShoot(suite.ctx, nonExistentID)

		// Then
		suite.Error(err, "should return error for non-existent shoot")
	})

}

func (suite *ShootRepositoryTestSuit) TestGetShoots() {
	// Given
	testShoot1 := testutils.CreateTestShoot(suite.testClient.Id)
	testShoot2 := testutils.CreateTestShootWithOptions(suite.testClient.Id, func(shoot *model.Shoot) {
		shoot.ShootDate = time.Now().AddDate(0, 0, 10)
		shoot.ShootLocation = "photo studio Aurora"
		shoot.Notes = ""
	})
	testShoot3 := testutils.CreateTestShootWithOptions(suite.testClient.Id, func(shoot *model.Shoot) {
		shoot.ShootDate = time.Now().AddDate(0, 0, 60)
		shoot.ShootLocation = "beacon"
		shoot.Notes = "blanket"
		shoot.ShootType = "family"
	})
	err := suite.repo.AddShoot(suite.ctx, testShoot1)
	suite.NoError(err, "shoot should be created")
	err = suite.repo.AddShoot(suite.ctx, testShoot2)
	suite.NoError(err, "shoot should be created")
	err = suite.repo.AddShoot(suite.ctx, testShoot3)
	suite.NoError(err, "shoot should be created")

	// When
	allShoots, err := suite.repo.GetShoots(suite.ctx)

	// Then
	suite.NoError(err, "should not return error")
	suite.Len(allShoots, 3, "should return 3 shoots")

	foundShoots := make(map[int]model.Shoot)
	for _, shoot := range allShoots {
		foundShoots[shoot.Id] = shoot
	}

	suite.Contains(foundShoots, testShoot1.Id, "should contain first shoot")
	suite.Contains(foundShoots, testShoot2.Id, "should contain second shoot")
	suite.Contains(foundShoots, testShoot3.Id, "should contain third shoot")

	shoot1 := foundShoots[testShoot1.Id]
	suite.Equal(testShoot1.ShootDate.Format("2006-01-02"), shoot1.ShootDate.Format("2006-01-02"),
		"shoot date should match")
	suite.Equal(testShoot1.StartTime.Format("15:04"), shoot1.StartTime.Format("15:04"),
		"shoot start time should match")
	suite.Equal(testShoot1.EndTime.Format("15:04"), shoot1.EndTime.Format("15:04"),
		"shoot end time should match")
	suite.Equal(testShoot1.ShootPrice, shoot1.ShootPrice, "shoot price should match")
	suite.Equal(testShoot1.ShootLocation, shoot1.ShootLocation, "shoot location should match")
	suite.Equal(testShoot1.ShootType, shoot1.ShootType, "shoot type should match")
	suite.Equal(testShoot1.Notes, shoot1.Notes, "notes should match")

	shoot2 := foundShoots[testShoot2.Id]
	suite.Equal(testShoot2.ShootDate.Format("2006-01-02"), shoot2.ShootDate.Format("2006-01-02"),
		"shoot date should match")
	suite.Equal(testShoot2.StartTime.Format("15:04"), shoot2.StartTime.Format("15:04"),
		"shoot start time should match")
	suite.Equal(testShoot2.EndTime.Format("15:04"), shoot2.EndTime.Format("15:04"),
		"shoot end time should match")
	suite.Equal(testShoot2.ShootPrice, shoot2.ShootPrice, "shoot price should match")
	suite.Equal(testShoot2.ShootLocation, shoot2.ShootLocation, "shoot location should match")
	suite.Equal(testShoot2.ShootType, shoot2.ShootType, "shoot type should match")
	suite.Equal(testShoot2.Notes, shoot2.Notes, "notes should match")

	shoot3 := foundShoots[testShoot3.Id]
	suite.Equal(testShoot3.ShootDate.Format("2006-01-02"), shoot3.ShootDate.Format("2006-01-02"),
		"shoot date should match")
	suite.Equal(testShoot3.StartTime.Format("15:04"), shoot3.StartTime.Format("15:04"),
		"shoot start time should match")
	suite.Equal(testShoot3.EndTime.Format("15:04"), shoot3.EndTime.Format("15:04"),
		"shoot end time should match")
	suite.Equal(testShoot3.ShootPrice, shoot3.ShootPrice, "shoot price should match")
	suite.Equal(testShoot3.ShootLocation, shoot3.ShootLocation, "shoot location should match")
	suite.Equal(testShoot3.ShootType, shoot3.ShootType, "shoot type should match")
	suite.Equal(testShoot3.Notes, shoot3.Notes, "notes should match")

}
