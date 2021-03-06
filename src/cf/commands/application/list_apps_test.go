package application_test

import (
	"cf"
	. "cf/commands/application"
	"cf/configuration"
	"github.com/stretchr/testify/assert"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
	"testing"
)

func TestApps(t *testing.T) {
	domain := cf.DomainFields{}
	domain.Name = "cfapps.io"
	domain2 := cf.DomainFields{}
	domain2.Name = "example.com"

	route1 := cf.RouteSummary{}
	route1.Host = "app1"
	route1.Domain = domain

	route2 := cf.RouteSummary{}
	route2.Host = "app1"
	route2.Domain = domain2

	app1Routes := []cf.RouteSummary{route1, route2}

	domain3 := cf.DomainFields{}
	domain3.Name = "cfapps.io"

	route3 := cf.RouteSummary{}
	route3.Host = "app2"
	route3.Domain = domain3

	app2Routes := []cf.RouteSummary{route3}

	app := cf.AppSummary{}
	app.Name = "Application-1"
	app.State = "started"
	app.RunningInstances = 1
	app.InstanceCount = 1
	app.Memory = 512
	app.DiskQuota = 1024
	app.RouteSummaries = app1Routes

	app2 := cf.AppSummary{}
	app2.Name = "Application-2"
	app2.State = "started"
	app2.RunningInstances = 1
	app2.InstanceCount = 2
	app2.Memory = 256
	app2.DiskQuota = 1024
	app2.RouteSummaries = app2Routes

	apps := []cf.AppSummary{app, app2}

	appSummaryRepo := &testapi.FakeAppSummaryRepo{
		GetSummariesInCurrentSpaceApps: apps,
	}

	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true}

	ui := callApps(t, appSummaryRepo, reqFactory)

	assert.True(t, testcmd.CommandDidPassRequirements)

	testassert.SliceContains(t, ui.Outputs, testassert.Lines{
		{"Getting apps in", "my-org", "development", "my-user"},
		{"OK"},
		{"Application-1", "started", "1/1", "512M", "1G", "app1.cfapps.io", "app1.example.com"},
		{"Application-2", "started", "1/2", "256M", "1G", "app2.cfapps.io"},
	})
}

func TestAppsEmptyList(t *testing.T) {
	appSummaryRepo := &testapi.FakeAppSummaryRepo{
		GetSummariesInCurrentSpaceApps: []cf.AppSummary{},
	}

	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: true}

	ui := callApps(t, appSummaryRepo, reqFactory)

	assert.True(t, testcmd.CommandDidPassRequirements)
	testassert.SliceContains(t, ui.Outputs, testassert.Lines{
		{"Getting apps in", "my-org", "development", "my-user"},
		{"OK"},
		{"No apps found"},
	})
}

func TestAppsRequiresLogin(t *testing.T) {
	appSummaryRepo := &testapi.FakeAppSummaryRepo{}
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: false, TargetedSpaceSuccess: true}

	callApps(t, appSummaryRepo, reqFactory)

	assert.False(t, testcmd.CommandDidPassRequirements)
}

func TestAppsRequiresASelectedSpaceAndOrg(t *testing.T) {
	appSummaryRepo := &testapi.FakeAppSummaryRepo{}
	reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedSpaceSuccess: false}

	callApps(t, appSummaryRepo, reqFactory)

	assert.False(t, testcmd.CommandDidPassRequirements)
}

func callApps(t *testing.T, appSummaryRepo *testapi.FakeAppSummaryRepo, reqFactory *testreq.FakeReqFactory) (ui *testterm.FakeUI) {
	ui = &testterm.FakeUI{}
	token, err := testconfig.CreateAccessTokenWithTokenInfo(configuration.TokenInfo{
		Username: "my-user",
	})
	assert.NoError(t, err)

	space := cf.SpaceFields{}
	space.Name = "development"
	org := cf.OrganizationFields{}
	org.Name = "my-org"
	config := &configuration.Configuration{
		SpaceFields:        space,
		OrganizationFields: org,
		AccessToken:        token,
	}

	ctxt := testcmd.NewContext("apps", []string{})
	cmd := NewListApps(ui, config, appSummaryRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)

	return
}
