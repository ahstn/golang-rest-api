package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testJets(t *testing.T) {
	t.Parallel()

	query := Jets(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testJetsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = jet.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Jets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testJetsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Jets(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Jets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testJetsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := JetSlice{jet}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Jets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testJetsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := JetExists(tx, jet.ID)
	if err != nil {
		t.Errorf("Unable to check if Jet exists: %s", err)
	}
	if !e {
		t.Errorf("Expected JetExistsG to return true, but got false.")
	}
}
func testJetsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	jetFound, err := FindJet(tx, jet.ID)
	if err != nil {
		t.Error(err)
	}

	if jetFound == nil {
		t.Error("want a record, got nil")
	}
}
func testJetsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Jets(tx).Bind(jet); err != nil {
		t.Error(err)
	}
}

func testJetsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Jets(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testJetsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jetOne := &Jet{}
	jetTwo := &Jet{}
	if err = randomize.Struct(seed, jetOne, jetDBTypes, false, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}
	if err = randomize.Struct(seed, jetTwo, jetDBTypes, false, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jetOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = jetTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Jets(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testJetsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	jetOne := &Jet{}
	jetTwo := &Jet{}
	if err = randomize.Struct(seed, jetOne, jetDBTypes, false, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}
	if err = randomize.Struct(seed, jetTwo, jetDBTypes, false, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jetOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = jetTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Jets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
func jetBeforeInsertHook(e boil.Executor, o *Jet) error {
	*o = Jet{}
	return nil
}

func jetAfterInsertHook(e boil.Executor, o *Jet) error {
	*o = Jet{}
	return nil
}

func jetAfterSelectHook(e boil.Executor, o *Jet) error {
	*o = Jet{}
	return nil
}

func jetBeforeUpdateHook(e boil.Executor, o *Jet) error {
	*o = Jet{}
	return nil
}

func jetAfterUpdateHook(e boil.Executor, o *Jet) error {
	*o = Jet{}
	return nil
}

func jetBeforeDeleteHook(e boil.Executor, o *Jet) error {
	*o = Jet{}
	return nil
}

func jetAfterDeleteHook(e boil.Executor, o *Jet) error {
	*o = Jet{}
	return nil
}

func jetBeforeUpsertHook(e boil.Executor, o *Jet) error {
	*o = Jet{}
	return nil
}

func jetAfterUpsertHook(e boil.Executor, o *Jet) error {
	*o = Jet{}
	return nil
}

func testJetsHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Jet{}
	o := &Jet{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, jetDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Jet object: %s", err)
	}

	AddJetHook(boil.BeforeInsertHook, jetBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	jetBeforeInsertHooks = []JetHook{}

	AddJetHook(boil.AfterInsertHook, jetAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	jetAfterInsertHooks = []JetHook{}

	AddJetHook(boil.AfterSelectHook, jetAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	jetAfterSelectHooks = []JetHook{}

	AddJetHook(boil.BeforeUpdateHook, jetBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	jetBeforeUpdateHooks = []JetHook{}

	AddJetHook(boil.AfterUpdateHook, jetAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	jetAfterUpdateHooks = []JetHook{}

	AddJetHook(boil.BeforeDeleteHook, jetBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	jetBeforeDeleteHooks = []JetHook{}

	AddJetHook(boil.AfterDeleteHook, jetAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	jetAfterDeleteHooks = []JetHook{}

	AddJetHook(boil.BeforeUpsertHook, jetBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	jetBeforeUpsertHooks = []JetHook{}

	AddJetHook(boil.AfterUpsertHook, jetAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	jetAfterUpsertHooks = []JetHook{}
}
func testJetsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Jets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testJetsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx, jetColumns...); err != nil {
		t.Error(err)
	}

	count, err := Jets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testJetToOnePilotUsingPilot(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Jet
	var foreign Pilot

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, jetDBTypes, true, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, pilotDBTypes, true, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.PilotID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Pilot(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := JetSlice{&local}
	if err = local.L.LoadPilot(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Pilot == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Pilot = nil
	if err = local.L.LoadPilot(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Pilot == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testJetToOneSetOpPilotUsingPilot(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Jet
	var b, c Pilot

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, jetDBTypes, false, strmangle.SetComplement(jetPrimaryKeyColumns, jetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, pilotDBTypes, false, strmangle.SetComplement(pilotPrimaryKeyColumns, pilotColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, pilotDBTypes, false, strmangle.SetComplement(pilotPrimaryKeyColumns, pilotColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Pilot{&b, &c} {
		err = a.SetPilot(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Pilot != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Jets[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.PilotID != x.ID {
			t.Error("foreign key was wrong value", a.PilotID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.PilotID))
		reflect.Indirect(reflect.ValueOf(&a.PilotID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.PilotID != x.ID {
			t.Error("foreign key was wrong value", a.PilotID, x.ID)
		}
	}
}
func testJetsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = jet.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testJetsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := JetSlice{jet}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testJetsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Jets(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	jetDBTypes = map[string]string{`Age`: `integer`, `Color`: `text`, `ID`: `integer`, `Name`: `text`, `PilotID`: `integer`}
	_          = bytes.MinRead
)

func testJetsUpdate(t *testing.T) {
	t.Parallel()

	if len(jetColumns) == len(jetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Jets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, jet, jetDBTypes, true, jetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	if err = jet.Update(tx); err != nil {
		t.Error(err)
	}
}

func testJetsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(jetColumns) == len(jetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	jet := &Jet{}
	if err = randomize.Struct(seed, jet, jetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Jets(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, jet, jetDBTypes, true, jetPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(jetColumns, jetPrimaryKeyColumns) {
		fields = jetColumns
	} else {
		fields = strmangle.SetComplement(
			jetColumns,
			jetPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(jet))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := JetSlice{jet}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testJetsUpsert(t *testing.T) {
	t.Parallel()

	if len(jetColumns) == len(jetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	jet := Jet{}
	if err = randomize.Struct(seed, &jet, jetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = jet.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Jet: %s", err)
	}

	count, err := Jets(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &jet, jetDBTypes, false, jetPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Jet struct: %s", err)
	}

	if err = jet.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Jet: %s", err)
	}

	count, err = Jets(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
