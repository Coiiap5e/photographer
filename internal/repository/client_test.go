package repository

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/Coiiap5e/photographer/testutils"
	"github.com/stretchr/testify/suite"
)

func TestClientRepositoryTestSuit(t *testing.T) {
	suite.Run(t, new(ClientRepositoryTestSuit))
}

type ClientRepositoryTestSuit struct {
	suite.Suite
	ctx    context.Context
	db     *testutils.TestDB
	repo   Client
	logger *slog.Logger
}

func (suite *ClientRepositoryTestSuit) SetupSuite() {
	suite.ctx = context.Background()

	suite.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var err error

	suite.db, err = testutils.CreateTestDB(suite.ctx, suite.logger)
	suite.Require().NoError(err, "Failed to setup test database")
	suite.Require().NotNil(suite.db, "TestDB should not be nil")
	suite.Require().NotNil(suite.db.GetDB(), "DB connection should not be nil")

	suite.repo = NewClient(suite.db.GetDB())
}

func (suite *ClientRepositoryTestSuit) SetupTest() {
	err := suite.db.CleanTables(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *ClientRepositoryTestSuit) TearDownSuite() {
	if suite.db != nil {
		if err := suite.db.Cleanup(suite.ctx); err != nil {
			suite.logger.Error("failed to cleanup test database", "error", err)
			os.Exit(1)
		}
	}
}

func (suite *ClientRepositoryTestSuit) TestAddClient() {
	// Given
	testClient := testutils.CreateTestClient()

	// When
	err := suite.repo.AddClient(suite.ctx, testClient)

	// Then
	suite.NoError(err, "should not return error")
	suite.Greater(testClient.Id, 0, "should assign positive ID to client object")
	suite.NotZero(testClient.CreatedAt, "should set creation timestamp in client object")
}

func (suite *ClientRepositoryTestSuit) TestGetClientByID() {
	suite.T().Run("should successfully get a client by ID", func(t *testing.T) {
		// Given
		testClient := testutils.CreateTestClient()
		err := suite.repo.AddClient(suite.ctx, testClient)
		suite.Require().NoError(err, "precondition: client should be created")

		// When
		retrievedClient, err := suite.repo.GetClientByID(suite.ctx, testClient.Id)

		// Then
		suite.NoError(err, "should retrieve client without error")
		suite.NotNil(retrievedClient, "should return client object")
		suite.Equal(testClient.Id, retrievedClient.Id, "should retrieve correct client by ID")
		suite.Equal(testClient.FirstName, retrievedClient.FirstName, "should persist first name")
		suite.Equal(testClient.LastName, retrievedClient.LastName, "should persist last name")
		suite.Equal(testClient.Phone, retrievedClient.Phone, "should persist phone")
		suite.Equal(testClient.SocialNetworkUrl, retrievedClient.SocialNetworkUrl, "should persist social network URL")
	})
	suite.T().Run("should return error when get non-existent client", func(t *testing.T) {
		// Given
		nonExistentID := 9999

		// When
		retrievedClient, err := suite.repo.GetClientByID(suite.ctx, nonExistentID)

		// Then
		suite.Error(err, "should return error for non-existent client")
		suite.Nil(retrievedClient, "should not return client object")
	})

}

func (suite *ClientRepositoryTestSuit) TestDeleteClient() {
	suite.T().Run("should successfully delete client", func(t *testing.T) {
		// Given
		testClient := testutils.CreateTestClient()
		err := suite.repo.AddClient(suite.ctx, testClient)
		suite.Require().NoError(err, "precondition: client should be created")

		clientID := testClient.Id

		// When
		err = suite.repo.DeleteClient(suite.ctx, testClient.Id)

		// Then
		suite.NoError(err, "should not return error")

		deletedClient, err := suite.repo.GetClientByID(suite.ctx, clientID)
		suite.Error(err, "should return error when getting deleted client")
		suite.Nil(deletedClient, "should not return deleted client")
	})
	suite.T().Run("should return error when deleting non-existent client", func(t *testing.T) {
		// Given
		nonExistentID := 9999

		// When
		err := suite.repo.DeleteClient(suite.ctx, nonExistentID)

		// Then
		suite.Error(err, "should return error for non-existent client")
	})

}

func (suite *ClientRepositoryTestSuit) TestGetClients() {
	// Given
	testClient1 := testutils.CreateTestClient()
	testClient2 := testutils.CreateTestClientWithOptions(func(client *model.Client) {
		client.FirstName = "Anna"
		client.LastName = "Petrova"
		client.Phone = "+7(901)111-11-11"
		client.SocialNetworkUrl = "@a.petrova"
	})
	testClient3 := testutils.CreateTestClientWithOptions(func(client *model.Client) {
		client.LastName = "Garkach"
		client.Phone = "+7(905)555-11-11"
	})
	err := suite.repo.AddClient(suite.ctx, testClient1)
	suite.NoError(err, "should not return error")
	err = suite.repo.AddClient(suite.ctx, testClient2)
	suite.NoError(err, "should not return error")
	err = suite.repo.AddClient(suite.ctx, testClient3)
	suite.NoError(err, "should not return error")

	// When
	allClients, err := suite.repo.GetClients(suite.ctx)

	// Then
	suite.NoError(err, "should not return error")
	suite.Len(allClients, 3, "should return 3 clients")

	foundClients := make(map[int]model.Client)
	for _, client := range allClients {
		foundClients[client.Id] = client
	}

	suite.Contains(foundClients, testClient1.Id, "should contain first client")
	suite.Contains(foundClients, testClient2.Id, "should contain second client")
	suite.Contains(foundClients, testClient3.Id, "should contain third client")

	client1 := foundClients[testClient1.Id]
	suite.Equal(testClient1.FirstName, client1.FirstName, "client1 first name should match")
	suite.Equal(testClient1.LastName, client1.LastName, "client1 last name should match")
	suite.Equal(testClient1.Phone, client1.Phone, "client1 phone should match")
	suite.Equal(testClient1.SocialNetworkUrl, client1.SocialNetworkUrl, "client1 social network url should match")

	client2 := foundClients[testClient2.Id]
	suite.Equal(testClient2.FirstName, client2.FirstName, "client2 first name should match")
	suite.Equal(testClient2.LastName, client2.LastName, "client2 last name should match")
	suite.Equal(testClient2.Phone, client2.Phone, "client2 phone should match")
	suite.Equal(testClient2.SocialNetworkUrl, client2.SocialNetworkUrl, "client2 social network url should match")

	client3 := foundClients[testClient3.Id]
	suite.Equal(testClient3.FirstName, client3.FirstName, "client3 first name should match")
	suite.Equal(testClient3.LastName, client3.LastName, "client3 last name should match")
	suite.Equal(testClient3.Phone, client3.Phone, "client3 phone should match")
	suite.Equal(testClient3.SocialNetworkUrl, client3.SocialNetworkUrl, "client3 social network url should match")

}
