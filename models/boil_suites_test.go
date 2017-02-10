package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Pilots", testPilots)
	t.Run("Jets", testJets)
}

func TestDelete(t *testing.T) {
	t.Run("Pilots", testPilotsDelete)
	t.Run("Jets", testJetsDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Pilots", testPilotsQueryDeleteAll)
	t.Run("Jets", testJetsQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Pilots", testPilotsSliceDeleteAll)
	t.Run("Jets", testJetsSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Pilots", testPilotsExists)
	t.Run("Jets", testJetsExists)
}

func TestFind(t *testing.T) {
	t.Run("Pilots", testPilotsFind)
	t.Run("Jets", testJetsFind)
}

func TestBind(t *testing.T) {
	t.Run("Pilots", testPilotsBind)
	t.Run("Jets", testJetsBind)
}

func TestOne(t *testing.T) {
	t.Run("Pilots", testPilotsOne)
	t.Run("Jets", testJetsOne)
}

func TestAll(t *testing.T) {
	t.Run("Pilots", testPilotsAll)
	t.Run("Jets", testJetsAll)
}

func TestCount(t *testing.T) {
	t.Run("Pilots", testPilotsCount)
	t.Run("Jets", testJetsCount)
}

func TestHooks(t *testing.T) {
	t.Run("Pilots", testPilotsHooks)
	t.Run("Jets", testJetsHooks)
}

func TestInsert(t *testing.T) {
	t.Run("Pilots", testPilotsInsert)
	t.Run("Pilots", testPilotsInsertWhitelist)
	t.Run("Jets", testJetsInsert)
	t.Run("Jets", testJetsInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("JetToPilotUsingPilot", testJetToOnePilotUsingPilot)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("PilotToJets", testPilotToManyJets)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("JetToPilotUsingPilot", testJetToOneSetOpPilotUsingPilot)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("PilotToJets", testPilotToManyAddOpJets)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {}

func TestReload(t *testing.T) {
	t.Run("Pilots", testPilotsReload)
	t.Run("Jets", testJetsReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Pilots", testPilotsReloadAll)
	t.Run("Jets", testJetsReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Pilots", testPilotsSelect)
	t.Run("Jets", testJetsSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Pilots", testPilotsUpdate)
	t.Run("Jets", testJetsUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Pilots", testPilotsSliceUpdateAll)
	t.Run("Jets", testJetsSliceUpdateAll)
}

func TestUpsert(t *testing.T) {
	t.Run("Pilots", testPilotsUpsert)
	t.Run("Jets", testJetsUpsert)
}
