package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testPilots(t *testing.T) {
	t.Parallel()

	query := Pilots(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testPilotsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = pilot.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Pilots(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testPilotsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Pilots(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Pilots(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testPilotsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := PilotSlice{pilot}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Pilots(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testPilotsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := PilotExists(tx, pilot.ID)
	if err != nil {
		t.Errorf("Unable to check if Pilot exists: %s", err)
	}
	if !e {
		t.Errorf("Expected PilotExistsG to return true, but got false.")
	}
}
func testPilotsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	pilotFound, err := FindPilot(tx, pilot.ID)
	if err != nil {
		t.Error(err)
	}

	if pilotFound == nil {
		t.Error("want a record, got nil")
	}
}
func testPilotsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Pilots(tx).Bind(pilot); err != nil {
		t.Error(err)
	}
}

func testPilotsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Pilots(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testPilotsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilotOne := &Pilot{}
	pilotTwo := &Pilot{}
	if err = randomize.Struct(seed, pilotOne, pilotDBTypes, false, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}
	if err = randomize.Struct(seed, pilotTwo, pilotDBTypes, false, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilotOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = pilotTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Pilots(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testPilotsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	pilotOne := &Pilot{}
	pilotTwo := &Pilot{}
	if err = randomize.Struct(seed, pilotOne, pilotDBTypes, false, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}
	if err = randomize.Struct(seed, pilotTwo, pilotDBTypes, false, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilotOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = pilotTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Pilots(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
func pilotBeforeInsertHook(e boil.Executor, o *Pilot) error {
	*o = Pilot{}
	return nil
}

func pilotAfterInsertHook(e boil.Executor, o *Pilot) error {
	*o = Pilot{}
	return nil
}

func pilotAfterSelectHook(e boil.Executor, o *Pilot) error {
	*o = Pilot{}
	return nil
}

func pilotBeforeUpdateHook(e boil.Executor, o *Pilot) error {
	*o = Pilot{}
	return nil
}

func pilotAfterUpdateHook(e boil.Executor, o *Pilot) error {
	*o = Pilot{}
	return nil
}

func pilotBeforeDeleteHook(e boil.Executor, o *Pilot) error {
	*o = Pilot{}
	return nil
}

func pilotAfterDeleteHook(e boil.Executor, o *Pilot) error {
	*o = Pilot{}
	return nil
}

func pilotBeforeUpsertHook(e boil.Executor, o *Pilot) error {
	*o = Pilot{}
	return nil
}

func pilotAfterUpsertHook(e boil.Executor, o *Pilot) error {
	*o = Pilot{}
	return nil
}

func testPilotsHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Pilot{}
	o := &Pilot{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, pilotDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Pilot object: %s", err)
	}

	AddPilotHook(boil.BeforeInsertHook, pilotBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	pilotBeforeInsertHooks = []PilotHook{}

	AddPilotHook(boil.AfterInsertHook, pilotAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	pilotAfterInsertHooks = []PilotHook{}

	AddPilotHook(boil.AfterSelectHook, pilotAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	pilotAfterSelectHooks = []PilotHook{}

	AddPilotHook(boil.BeforeUpdateHook, pilotBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	pilotBeforeUpdateHooks = []PilotHook{}

	AddPilotHook(boil.AfterUpdateHook, pilotAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	pilotAfterUpdateHooks = []PilotHook{}

	AddPilotHook(boil.BeforeDeleteHook, pilotBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	pilotBeforeDeleteHooks = []PilotHook{}

	AddPilotHook(boil.AfterDeleteHook, pilotAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	pilotAfterDeleteHooks = []PilotHook{}

	AddPilotHook(boil.BeforeUpsertHook, pilotBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	pilotBeforeUpsertHooks = []PilotHook{}

	AddPilotHook(boil.AfterUpsertHook, pilotAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	pilotAfterUpsertHooks = []PilotHook{}
}
func testPilotsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Pilots(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testPilotsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx, pilotColumns...); err != nil {
		t.Error(err)
	}

	count, err := Pilots(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testPilotToManyJets(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Pilot
	var b, c Jet

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, pilotDBTypes, true, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, jetDBTypes, false, jetColumnsWithDefault...)
	randomize.Struct(seed, &c, jetDBTypes, false, jetColumnsWithDefault...)

	b.PilotID = a.ID
	c.PilotID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	jet, err := a.Jets(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range jet {
		if v.PilotID == b.PilotID {
			bFound = true
		}
		if v.PilotID == c.PilotID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := PilotSlice{&a}
	if err = a.L.LoadJets(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Jets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Jets = nil
	if err = a.L.LoadJets(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Jets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", jet)
	}
}

func testPilotToManyAddOpJets(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Pilot
	var b, c, d, e Jet

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, pilotDBTypes, false, strmangle.SetComplement(pilotPrimaryKeyColumns, pilotColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Jet{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, jetDBTypes, false, strmangle.SetComplement(jetPrimaryKeyColumns, jetColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Jet{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddJets(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.PilotID {
			t.Error("foreign key was wrong value", a.ID, first.PilotID)
		}
		if a.ID != second.PilotID {
			t.Error("foreign key was wrong value", a.ID, second.PilotID)
		}

		if first.R.Pilot != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Pilot != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Jets[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Jets[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Jets(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testPilotsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = pilot.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testPilotsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := PilotSlice{pilot}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testPilotsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Pilots(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	pilotDBTypes = map[string]string{`ID`: `integer`, `Name`: `text`}
	_            = bytes.MinRead
)

func testPilotsUpdate(t *testing.T) {
	t.Parallel()

	if len(pilotColumns) == len(pilotPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Pilots(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, pilot, pilotDBTypes, true, pilotColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	if err = pilot.Update(tx); err != nil {
		t.Error(err)
	}
}

func testPilotsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(pilotColumns) == len(pilotPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	pilot := &Pilot{}
	if err = randomize.Struct(seed, pilot, pilotDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Pilots(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, pilot, pilotDBTypes, true, pilotPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(pilotColumns, pilotPrimaryKeyColumns) {
		fields = pilotColumns
	} else {
		fields = strmangle.SetComplement(
			pilotColumns,
			pilotPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(pilot))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := PilotSlice{pilot}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testPilotsUpsert(t *testing.T) {
	t.Parallel()

	if len(pilotColumns) == len(pilotPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	pilot := Pilot{}
	if err = randomize.Struct(seed, &pilot, pilotDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = pilot.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Pilot: %s", err)
	}

	count, err := Pilots(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &pilot, pilotDBTypes, false, pilotPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Pilot struct: %s", err)
	}

	if err = pilot.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Pilot: %s", err)
	}

	count, err = Pilots(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
