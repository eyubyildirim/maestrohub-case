package game

import (
	"time"

	"github.com/google/uuid"
)

type UnitType string

const (
	Soldier    UnitType = "soldier"
	Tank       UnitType = "tank"
	Helicopter UnitType = "helicopter"
)

type Tower struct {
	ID string `json:"id"`
}

type Unit struct {
	ID            string   `json:"id"`
	Type          UnitType `json:"type"`
	TTL           float64  `json:"ttl"`
	CreationTime  int64    `json:"creationTime"`
	TransitTTL    float64  `json:"transitTtl"`
	MaxTransitTTL float64  `json:"maxTransitTtl"`
	HP            int      `json:"hp"`
	MaxHP         int      `json:"maxHp"`
}

type Group struct {
	ID            string  `json:"id"`
	Units         []*Unit `json:"units"`
	TransitTTL    float64 `json:"transitTtl"`
	MaxTransitTTL float64 `json:"maxTransitTtl"`
}

type Squad struct {
	ID          string           `json:"id"`
	Composition map[UnitType]int `json:"composition"`
}

type GameState struct {
	AvailableUnits           map[string]*Unit  `json:"availableUnits"`
	IndividualsInTransit     map[string]*Unit  `json:"individualsInTransit"`
	GroupsInTransit          map[string]*Group `json:"groupsInTransit"`
	Score                    int               `json:"score"`
	BreachCount              int               `json:"breachCount"`
	EscapedCount             int               `json:"escapedCount"`
	LastEvents               []string          `json:"lastEvents"`
	Squads                   map[string]*Squad `json:"squads"`
	ProcessingIndividualUnit *Unit             `json:"processingIndividualUnit"`
	ProcessingGroup          *Group            `json:"processingGroup"`
}

func NewGameState() *GameState {
	return &GameState{
		AvailableUnits:       make(map[string]*Unit),
		IndividualsInTransit: make(map[string]*Unit),
		GroupsInTransit:      make(map[string]*Group),
		Squads:               make(map[string]*Squad),
		LastEvents:           make([]string, 0),
	}
}

func NewSquad(composition map[UnitType]int) *Squad {
	return &Squad{
		ID:          uuid.NewString(),
		Composition: composition,
	}
}

func NewUnit(unitType UnitType, initialTTL float64, initialHP int) *Unit {
	return &Unit{
		ID:           uuid.NewString(),
		Type:         unitType,
		TTL:          initialTTL,
		CreationTime: time.Now().Unix(),
		HP:           initialHP,
		MaxHP:        initialHP,
	}
}

func NewGroup(units []*Unit) *Group {
	return &Group{
		ID:    uuid.NewString(),
		Units: units,
	}
}
