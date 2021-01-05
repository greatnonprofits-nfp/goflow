package test

import (
	"io/ioutil"

	"github.com/nyaruka/gocommon/uuids"
	"github.com/greatnonprofits-nfp/goflow/assets"
	"github.com/greatnonprofits-nfp/goflow/assets/static"
	"github.com/greatnonprofits-nfp/goflow/assets/static/types"
	"github.com/greatnonprofits-nfp/goflow/envs"
	"github.com/greatnonprofits-nfp/goflow/flows"
	"github.com/greatnonprofits-nfp/goflow/flows/definition/migrations"
	"github.com/greatnonprofits-nfp/goflow/flows/engine"
)

// LoadSessionAssets loads a session assets instance from a static JSON file
func LoadSessionAssets(env envs.Environment, path string) (flows.SessionAssets, error) {
	assetsJSON, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	source, err := static.NewSource(assetsJSON)
	if err != nil {
		return nil, err
	}

	mconfig := &migrations.Config{BaseMediaURL: "http://temba.io/"}

	return engine.NewSessionAssets(env, source, mconfig)
}

func LoadFlowFromAssets(env envs.Environment, path string, uuid assets.FlowUUID) (flows.Flow, error) {
	sa, err := LoadSessionAssets(env, path)
	if err != nil {
		return nil, err
	}

	return sa.Flows().Get(uuid)
}

func NewChannel(name string, address string, schemes []string, roles []assets.ChannelRole, parent *assets.ChannelReference) *flows.Channel {
	return flows.NewChannel(types.NewChannel(assets.ChannelUUID(uuids.New()), name, address, schemes, roles, parent))
}

func NewTelChannel(name string, address string, roles []assets.ChannelRole, parent *assets.ChannelReference, country envs.Country, matchPrefixes []string, allowInternational bool) *flows.Channel {
	return flows.NewChannel(types.NewTelChannel(assets.ChannelUUID(uuids.New()), name, address, roles, parent, country, matchPrefixes, allowInternational))
}

func NewClassifier(name, type_ string, intents []string) *flows.Classifier {
	return flows.NewClassifier(types.NewClassifier(assets.ClassifierUUID(uuids.New()), name, type_, intents))
}
