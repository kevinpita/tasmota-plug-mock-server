package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite

	server     *Server
	testServer *httptest.Server
}

func (suite *MainTestSuite) SetupSuite() {
	suite.server = NewServer()
	suite.testServer = httptest.NewServer(http.HandlerFunc(suite.server.handleCmnd))
}

func (suite *MainTestSuite) TearDownSuite() {
	suite.testServer.Close()
}

func (suite *MainTestSuite) TestPowerOn() {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		suite.testServer.URL+"/?cmnd=Power%20On",
		nil,
	)
	suite.Require().NoError(err)
	resp, err := suite.testServer.Client().Do(req)
	suite.Require().NoError(err)
	defer func() {
		errBodyClose := resp.Body.Close()
		suite.Require().NoError(errBodyClose)
	}()

	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var jsonResponse map[string]string
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)

	suite.Require().NoError(err)
	suite.Require().Equal("ON", jsonResponse["POWER"])
	suite.Require().Equal("ON", suite.server.powerState)
}

func (suite *MainTestSuite) TestPowerOff() {
	suite.server.powerState = "ON"

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		suite.testServer.URL+"/?cmnd=Power%20Off",
		nil,
	)
	suite.Require().NoError(err)
	resp, err := suite.testServer.Client().Do(req)
	suite.Require().NoError(err)
	defer func() {
		errBodyClose := resp.Body.Close()
		suite.Require().NoError(errBodyClose)
	}()

	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var jsonResponse map[string]string
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)

	suite.Require().NoError(err)
	suite.Require().Equal("OFF", jsonResponse["POWER"])
	suite.Require().Equal("OFF", suite.server.powerState)
}

func (suite *MainTestSuite) TestPowerStatus() {
	suite.server.powerState = "OFF"

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		suite.testServer.URL+"/?cmnd=Power",
		nil,
	)
	suite.Require().NoError(err)
	resp, err := suite.testServer.Client().Do(req)
	suite.Require().NoError(err)
	defer func() {
		errBodyClose := resp.Body.Close()
		suite.Require().NoError(errBodyClose)
	}()

	var jsonResponseOff map[string]string
	err = json.NewDecoder(resp.Body).Decode(&jsonResponseOff)

	suite.Require().NoError(err)
	suite.Require().Equal("OFF", jsonResponseOff["POWER"])

	suite.server.powerState = "ON"

	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		suite.testServer.URL+"/?cmnd=Power",
		nil,
	)
	suite.Require().NoError(err)
	resp, err = suite.testServer.Client().Do(req)
	suite.Require().NoError(err)
	defer func() {
		errBodyClose := resp.Body.Close()
		suite.Require().NoError(errBodyClose)
	}()

	var jsonResponseOn map[string]string
	err = json.NewDecoder(resp.Body).Decode(&jsonResponseOn)

	suite.Require().NoError(err)
	suite.Require().Equal("ON", jsonResponseOn["POWER"])
}

func (suite *MainTestSuite) TestEnergyTotal() {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		suite.testServer.URL+"/?cmnd=EnergyTotal",
		nil,
	)
	suite.Require().NoError(err)
	resp, err := suite.testServer.Client().Do(req)
	suite.Require().NoError(err)
	defer func() {
		errBodyClose := resp.Body.Close()
		suite.Require().NoError(errBodyClose)
	}()

	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var jsonResponse map[string]map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)

	suite.Require().NoError(err)
	suite.Require().Contains(jsonResponse, "EnergyTotal")
}

func (suite *MainTestSuite) TestNotFound() {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		suite.testServer.URL+"/?cmnd=UnknownCommand",
		nil,
	)
	suite.Require().NoError(err)
	resp, err := suite.testServer.Client().Do(req)
	suite.Require().NoError(err)
	defer func() {
		errBodyClose := resp.Body.Close()
		suite.Require().NoError(errBodyClose)
	}()

	suite.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
