package game

import (
	"log"
	"math/rand"
	"time"
	"tower-defence-engine/internal/config"
)

type Broadcaster interface {
	Broadcast(message any)
}

type ICommand any
type CreateSquadCommand struct{ UnitIDs []string }
type DeploySquadCommand struct{ SquadID string }
type DeployIndividualsCommand struct{ UnitIDs []string }

type Engine struct {
	GameState           *GameState
	broadcaster         Broadcaster
	Config              *config.Config
	Towers              []*Tower
	CommandCh           chan ICommand
	individualUnitQueue chan string
	groupQueue          chan *Group
}

func NewEngine(broadcaster Broadcaster, cfg *config.Config) *Engine {
	return &Engine{
		GameState:           NewGameState(),
		broadcaster:         broadcaster,
		Config:              cfg,
		Towers:              []*Tower{{ID: "T1"}, {ID: "T2"}, {ID: "T3"}},
		CommandCh:           make(chan ICommand, 128),
		individualUnitQueue: make(chan string, 256),
		groupQueue:          make(chan *Group, 64),
	}
}

func (e *Engine) Run() {
	spawnRate := time.Duration(e.Config.Rates.SpawnRateMs) * time.Millisecond
	unitWeaponRate := time.Duration(e.Config.Rates.IndividualWeaponRateMs) * time.Millisecond
	groupWeaponRate := time.Duration(e.Config.Rates.GroupWeaponRateMs) * time.Millisecond
	ttlCheckRate := time.Duration(e.Config.Rates.TtlCheckRateMs) * time.Millisecond

	spawnTicker := time.NewTicker(spawnRate)
	unitWeaponTicker := time.NewTicker(unitWeaponRate)
	groupWeaponTicker := time.NewTicker(groupWeaponRate)
	ttlTicker := time.NewTicker(ttlCheckRate)

	defer spawnTicker.Stop()
	defer unitWeaponTicker.Stop()
	defer groupWeaponTicker.Stop()
	defer ttlTicker.Stop()

	log.Println("Game Engine armed and running with HP mechanics.")

	for {
		var stateChanged bool

		select {
		case <-spawnTicker.C:
			e.spawnUnits()
			stateChanged = true
		case cmd := <-e.CommandCh:
			e.processCommand(cmd)
			stateChanged = true
		case <-ttlTicker.C:
			if e.updateAllTTLs() {
				stateChanged = true
			}
			if e.GameState.ProcessingIndividualUnit == nil {
				select {
				case unitID := <-e.individualUnitQueue:
					if unit, ok := e.GameState.IndividualsInTransit[unitID]; ok {
						e.GameState.ProcessingIndividualUnit = unit
						e.GameState.LastEvents = append(e.GameState.LastEvents, "LOCK_ON_INDIVIDUAL:"+unitID)
						log.Printf("🔫 Individual weapon locked on %s.", unitID)
						stateChanged = true
					}
				default:
				}
			}
			if e.GameState.ProcessingGroup == nil {
				select {
				case group := <-e.groupQueue:
					if _, ok := e.GameState.GroupsInTransit[group.ID]; ok {
						e.GameState.ProcessingGroup = group
						e.GameState.LastEvents = append(e.GameState.LastEvents, "LOCK_ON_GROUP:"+group.ID)
						log.Printf("🎯 Group weapon locked on %s.", group.ID)
						stateChanged = true
					}
				default:
				}
			}
		case <-unitWeaponTicker.C:
			if target := e.GameState.ProcessingIndividualUnit; target != nil {
				target.HP -= e.Config.Weapons.ProcessingPower.Individual
				log.Printf("Individual weapon applying damage to %s. HP left: %d", target.ID, target.HP)

				if target.HP <= 0 {
					e.GameState.Score += e.getPointsForUnit(target.Type)
					e.GameState.LastEvents = append(e.GameState.LastEvents, "DESTROY_INDIVIDUAL:"+target.ID)
					delete(e.GameState.IndividualsInTransit, target.ID)
					e.GameState.ProcessingIndividualUnit = nil
					log.Printf("💥 Individual weapon destroyed %s.", target.ID)
				}
				stateChanged = true
			}
		case <-groupWeaponTicker.C:
			if targetGroup := e.GameState.ProcessingGroup; targetGroup != nil {
				survivingUnits := make([]*Unit, 0)
				for _, unit := range targetGroup.Units {
					if unit.HP > 0 {
						unit.HP -= e.Config.Weapons.ProcessingPower.Group
						if unit.HP <= 0 {
							e.GameState.Score += e.getPointsForUnit(unit.Type)
						}
					}
					if unit.HP > 0 {
						survivingUnits = append(survivingUnits, unit)
					}
				}

				if len(survivingUnits) == 0 {
					e.GameState.LastEvents = append(e.GameState.LastEvents, "DESTROY_GROUP:"+targetGroup.ID)
					delete(e.GameState.GroupsInTransit, targetGroup.ID)
					e.GameState.ProcessingGroup = nil
					log.Printf("💣 Group %s completely destroyed.", targetGroup.ID)
				} else {
					targetGroup.Units = survivingUnits
				}
				stateChanged = true
			}
		}

		if stateChanged {
			e.broadcaster.Broadcast(e.GameState)
			if len(e.GameState.LastEvents) > 0 {
				e.GameState.LastEvents = make([]string, 0)
			}
		}
	}
}

func (e *Engine) spawnUnits() {
	unitTypes := []UnitType{Soldier, Tank, Helicopter}
	for range e.Towers {
		randomType := unitTypes[rand.Intn(len(unitTypes))]
		hp := e.getHpForUnit(randomType)
		newUnit := NewUnit(randomType, e.Config.Timing.BattlefieldTtlSec, hp)
		e.GameState.AvailableUnits[newUnit.ID] = newUnit
	}
}

func (e *Engine) processCommand(cmd ICommand) {
	transitTtl := e.Config.Timing.TransitTtlSec
	switch c := cmd.(type) {
	case CreateSquadCommand:
		composition := make(map[UnitType]int)
		for _, unitID := range c.UnitIDs {
			if unit, ok := e.GameState.AvailableUnits[unitID]; ok {
				composition[unit.Type]++
			}
		}
		if len(composition) > 0 {
			newSquad := NewSquad(composition)
			e.GameState.Squads[newSquad.ID] = newSquad
		}
	case DeployIndividualsCommand:
		for _, unitID := range c.UnitIDs {
			if unit, ok := e.GameState.AvailableUnits[unitID]; ok {
				delete(e.GameState.AvailableUnits, unitID)
				unit.TransitTTL, unit.MaxTransitTTL = transitTtl, transitTtl
				e.GameState.IndividualsInTransit[unit.ID] = unit
				e.individualUnitQueue <- unit.ID
			}
		}
	case DeploySquadCommand:
		squad, ok := e.GameState.Squads[c.SquadID]
		if !ok {
			return
		}
		availableCounts := make(map[UnitType]int)
		unitsByType := make(map[UnitType][]*Unit)
		for _, unit := range e.GameState.AvailableUnits {
			availableCounts[unit.Type]++
			unitsByType[unit.Type] = append(unitsByType[unit.Type], unit)
		}
		unitsToDeploy := make([]*Unit, 0)
		for unitType, neededCount := range squad.Composition {
			canDeployCount := min(neededCount, availableCounts[unitType])
			if canDeployCount > 0 {
				grabbedUnits := unitsByType[unitType][:canDeployCount]
				unitsToDeploy = append(unitsToDeploy, grabbedUnits...)
				for _, grabbedUnit := range grabbedUnits {
					delete(e.GameState.AvailableUnits, grabbedUnit.ID)
				}
			}
		}
		if len(unitsToDeploy) > 0 {
			newGroup := NewGroup(unitsToDeploy)
			newGroup.TransitTTL, newGroup.MaxTransitTTL = transitTtl, transitTtl
			e.GameState.GroupsInTransit[newGroup.ID] = newGroup
			e.groupQueue <- newGroup
		}
	}
}

func (e *Engine) updateAllTTLs() bool {
	tickDecrement := float64(e.Config.Rates.TtlCheckRateMs) / 1000.0
	escapedIDs := []string{}
	for _, unit := range e.GameState.AvailableUnits {
		unit.TTL -= tickDecrement
		if unit.TTL <= 0 {
			escapedIDs = append(escapedIDs, unit.ID)
		}
	}
	if len(escapedIDs) > 0 {
		for _, id := range escapedIDs {
			if _, ok := e.GameState.AvailableUnits[id]; ok {
				e.GameState.Score -= e.Config.Scoring.Penalties.Escape
				e.GameState.EscapedCount++
				delete(e.GameState.AvailableUnits, id)
				e.GameState.LastEvents = append(e.GameState.LastEvents, "ESCAPE:"+id)
			}
		}
	}
	breachedIndividualIDs := []string{}
	for _, unit := range e.GameState.IndividualsInTransit {
		unit.TransitTTL -= tickDecrement
		if unit.TransitTTL <= 0 {
			breachedIndividualIDs = append(breachedIndividualIDs, unit.ID)
		}
	}
	if len(breachedIndividualIDs) > 0 {
		for _, id := range breachedIndividualIDs {
			if e.GameState.ProcessingIndividualUnit == nil || e.GameState.ProcessingIndividualUnit.ID != id {
				delete(e.GameState.IndividualsInTransit, id)
				e.GameState.Score -= e.Config.Scoring.Penalties.Breach
				e.GameState.BreachCount++
				e.GameState.LastEvents = append(e.GameState.LastEvents, "BREACH_INDIVIDUAL:"+id)
			}
		}
	}
	breachedGroupIDs := []string{}
	for _, group := range e.GameState.GroupsInTransit {
		group.TransitTTL -= tickDecrement
		if group.TransitTTL <= 0 {
			breachedGroupIDs = append(breachedGroupIDs, group.ID)
		}
	}
	if len(breachedGroupIDs) > 0 {
		for _, id := range breachedGroupIDs {
			if e.GameState.ProcessingGroup == nil || e.GameState.ProcessingGroup.ID != id {
				delete(e.GameState.GroupsInTransit, id)
				e.GameState.Score -= e.Config.Scoring.Penalties.Breach
				e.GameState.BreachCount++
				e.GameState.LastEvents = append(e.GameState.LastEvents, "BREACH_GROUP:"+id)
			}
		}
	}
	return true
}

func (e *Engine) getPointsForUnit(unitType UnitType) int {
	switch unitType {
	case Soldier:
		return e.Config.Scoring.Points.Soldier
	case Tank:
		return e.Config.Scoring.Points.Tank
	case Helicopter:
		return e.Config.Scoring.Points.Helicopter
	default:
		return 0
	}
}

func (e *Engine) getHpForUnit(unitType UnitType) int {
	switch unitType {
	case Soldier:
		return e.Config.Units.Hitpoints.Soldier
	case Tank:
		return e.Config.Units.Hitpoints.Tank
	case Helicopter:
		return e.Config.Units.Hitpoints.Helicopter
	default:
		return 1
	}
}
