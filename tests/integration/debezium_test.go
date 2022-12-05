//go:build integ

package debezium_test

import (
	"github.com/dailydotdev/debezium_test/models"
	"github.com/dailydotdev/platform-go-common/ext/pgsql"
	"github.com/dailydotdev/platform-go-common/util/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"testing"
)

type debeziumTestSuite struct {
	suite.Suite
	pg *pgsql.Wrapper
}

func TestDebeziumTestSuite(t *testing.T) {
	suite.Run(t, new(debeziumTestSuite))
}

func (s *debeziumTestSuite) SetupSuite() {

	var err error

	s.pg, err = pgsql.NewWrapper(log.Logger.Level(zerolog.DebugLevel),
		pgsql.WithPrefix("debezium_test", "test"),
	)
	s.Require().NoError(err)
}

func (s *debeziumTestSuite) TearDownSuite() {
}

func (s *debeziumTestSuite) TestDebeziumTable1() {

	// will be processed via debezium-test-cdc-sub subscription

	err := s.pg.GetConnection().Create(&models.Table1{
		ID:   uuid.New(),
		Data: "1q2w3e",
	}).Error

	s.Require().NoError(err)
}

func (s *debeziumTestSuite) TestDebeziumTable2() {

	// will be processed via debezium-test-cdc-sub subscription

	err := s.pg.GetConnection().Create(&models.Table2{
		ID:   uuid.New(),
		Data: "1q2w3e",
	}).Error

	s.Require().NoError(err)
}

func (s *debeziumTestSuite) TestDebeziumTable3() {

	// will be processed via debezium-test-table3-sub subscription

	err := s.pg.GetConnection().Create(&models.Table3{
		ID:   uuid.New(),
		Data: "1q2w3e",
	}).Error

	s.Require().NoError(err)
}