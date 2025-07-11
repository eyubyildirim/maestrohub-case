export type UnitType = 'soldier' | 'tank' | 'helicopter'

export interface Unit {
  id: string
  type: UnitType
  ttl: number
  creationTime: number
  transitTtl: number
  maxTransitTtl: number
  hp: number
  maxHp: number
}

export interface Group {
  id: string
  units: Unit[]
  transitTtl: number
  maxTransitTtl: number
}

export interface Squad {
  id: string
  composition: Record<UnitType, number>
}

export interface GameState {
  availableUnits: Record<string, Unit>
  groupsInTransit: Record<string, Group>
  individualsInTransit: Record<string, Unit>
  score: number
  escapedCount: number
  breachCount: number
  lastEvents: string[]
  squads: Record<string, Squad>
  processingIndividualUnit: Unit | null
  processingGroup: Group | null
}

export interface Command {
  action: string
  payload: any
}

export interface CreateSquadPayload {
  unit_ids: string[]
}
export interface DeploySquadPayload {
  squad_id: string
}
export interface DeployIndividualsPayload {
  unit_ids: string[]
}
