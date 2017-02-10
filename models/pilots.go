package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/vattle/sqlboiler/strmangle"
)

// Pilot is an object representing the database table.
type Pilot struct {
	ID   int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name string `boil:"name" json:"name" toml:"name" yaml:"name"`

	R *pilotR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L pilotL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// pilotR is where relationships are stored.
type pilotR struct {
	Jets JetSlice
}

// pilotL is where Load methods for each relationship are stored.
type pilotL struct{}

var (
	pilotColumns               = []string{"id", "name"}
	pilotColumnsWithoutDefault = []string{"name"}
	pilotColumnsWithDefault    = []string{"id"}
	pilotPrimaryKeyColumns     = []string{"id"}
)

type (
	// PilotSlice is an alias for a slice of pointers to Pilot.
	// This should generally be used opposed to []Pilot.
	PilotSlice []*Pilot
	// PilotHook is the signature for custom Pilot hook methods
	PilotHook func(boil.Executor, *Pilot) error

	pilotQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	pilotType                 = reflect.TypeOf(&Pilot{})
	pilotMapping              = queries.MakeStructMapping(pilotType)
	pilotPrimaryKeyMapping, _ = queries.BindMapping(pilotType, pilotMapping, pilotPrimaryKeyColumns)
	pilotInsertCacheMut       sync.RWMutex
	pilotInsertCache          = make(map[string]insertCache)
	pilotUpdateCacheMut       sync.RWMutex
	pilotUpdateCache          = make(map[string]updateCache)
	pilotUpsertCacheMut       sync.RWMutex
	pilotUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var pilotBeforeInsertHooks []PilotHook
var pilotBeforeUpdateHooks []PilotHook
var pilotBeforeDeleteHooks []PilotHook
var pilotBeforeUpsertHooks []PilotHook

var pilotAfterInsertHooks []PilotHook
var pilotAfterSelectHooks []PilotHook
var pilotAfterUpdateHooks []PilotHook
var pilotAfterDeleteHooks []PilotHook
var pilotAfterUpsertHooks []PilotHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Pilot) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range pilotBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Pilot) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range pilotBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Pilot) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range pilotBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Pilot) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range pilotBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Pilot) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range pilotAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Pilot) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range pilotAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Pilot) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range pilotAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Pilot) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range pilotAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Pilot) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range pilotAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddPilotHook registers your hook function for all future operations.
func AddPilotHook(hookPoint boil.HookPoint, pilotHook PilotHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		pilotBeforeInsertHooks = append(pilotBeforeInsertHooks, pilotHook)
	case boil.BeforeUpdateHook:
		pilotBeforeUpdateHooks = append(pilotBeforeUpdateHooks, pilotHook)
	case boil.BeforeDeleteHook:
		pilotBeforeDeleteHooks = append(pilotBeforeDeleteHooks, pilotHook)
	case boil.BeforeUpsertHook:
		pilotBeforeUpsertHooks = append(pilotBeforeUpsertHooks, pilotHook)
	case boil.AfterInsertHook:
		pilotAfterInsertHooks = append(pilotAfterInsertHooks, pilotHook)
	case boil.AfterSelectHook:
		pilotAfterSelectHooks = append(pilotAfterSelectHooks, pilotHook)
	case boil.AfterUpdateHook:
		pilotAfterUpdateHooks = append(pilotAfterUpdateHooks, pilotHook)
	case boil.AfterDeleteHook:
		pilotAfterDeleteHooks = append(pilotAfterDeleteHooks, pilotHook)
	case boil.AfterUpsertHook:
		pilotAfterUpsertHooks = append(pilotAfterUpsertHooks, pilotHook)
	}
}

// OneP returns a single pilot record from the query, and panics on error.
func (q pilotQuery) OneP() *Pilot {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single pilot record from the query.
func (q pilotQuery) One() (*Pilot, error) {
	o := &Pilot{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for pilots")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Pilot records from the query, and panics on error.
func (q pilotQuery) AllP() PilotSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Pilot records from the query.
func (q pilotQuery) All() (PilotSlice, error) {
	var o PilotSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Pilot slice")
	}

	if len(pilotAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Pilot records in the query, and panics on error.
func (q pilotQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Pilot records in the query.
func (q pilotQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count pilots rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q pilotQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q pilotQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if pilots exists")
	}

	return count > 0, nil
}

// JetsG retrieves all the jet's jets.
func (o *Pilot) JetsG(mods ...qm.QueryMod) jetQuery {
	return o.Jets(boil.GetDB(), mods...)
}

// Jets retrieves all the jet's jets with an executor.
func (o *Pilot) Jets(exec boil.Executor, mods ...qm.QueryMod) jetQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"pilot_id\"=?", o.ID),
	)

	query := Jets(exec, queryMods...)
	queries.SetFrom(query.Query, "\"jets\" as \"a\"")
	return query
}

// LoadJets allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (pilotL) LoadJets(e boil.Executor, singular bool, maybePilot interface{}) error {
	var slice []*Pilot
	var object *Pilot

	count := 1
	if singular {
		object = maybePilot.(*Pilot)
	} else {
		slice = *maybePilot.(*PilotSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &pilotR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &pilotR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"jets\" where \"pilot_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load jets")
	}
	defer results.Close()

	var resultSlice []*Jet
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice jets")
	}

	if len(jetAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Jets = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.PilotID {
				local.R.Jets = append(local.R.Jets, foreign)
				break
			}
		}
	}

	return nil
}

// AddJetsG adds the given related objects to the existing relationships
// of the pilot, optionally inserting them as new records.
// Appends related to o.R.Jets.
// Sets related.R.Pilot appropriately.
// Uses the global database handle.
func (o *Pilot) AddJetsG(insert bool, related ...*Jet) error {
	return o.AddJets(boil.GetDB(), insert, related...)
}

// AddJetsP adds the given related objects to the existing relationships
// of the pilot, optionally inserting them as new records.
// Appends related to o.R.Jets.
// Sets related.R.Pilot appropriately.
// Panics on error.
func (o *Pilot) AddJetsP(exec boil.Executor, insert bool, related ...*Jet) {
	if err := o.AddJets(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddJetsGP adds the given related objects to the existing relationships
// of the pilot, optionally inserting them as new records.
// Appends related to o.R.Jets.
// Sets related.R.Pilot appropriately.
// Uses the global database handle and panics on error.
func (o *Pilot) AddJetsGP(insert bool, related ...*Jet) {
	if err := o.AddJets(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddJets adds the given related objects to the existing relationships
// of the pilot, optionally inserting them as new records.
// Appends related to o.R.Jets.
// Sets related.R.Pilot appropriately.
func (o *Pilot) AddJets(exec boil.Executor, insert bool, related ...*Jet) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.PilotID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"jets\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"pilot_id"}),
				strmangle.WhereClause("\"", "\"", 2, jetPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.PilotID = o.ID
		}
	}

	if o.R == nil {
		o.R = &pilotR{
			Jets: related,
		}
	} else {
		o.R.Jets = append(o.R.Jets, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &jetR{
				Pilot: o,
			}
		} else {
			rel.R.Pilot = o
		}
	}
	return nil
}

// PilotsG retrieves all records.
func PilotsG(mods ...qm.QueryMod) pilotQuery {
	return Pilots(boil.GetDB(), mods...)
}

// Pilots retrieves all the records using an executor.
func Pilots(exec boil.Executor, mods ...qm.QueryMod) pilotQuery {
	mods = append(mods, qm.From("\"pilots\""))
	return pilotQuery{NewQuery(exec, mods...)}
}

// FindPilotG retrieves a single record by ID.
func FindPilotG(id int, selectCols ...string) (*Pilot, error) {
	return FindPilot(boil.GetDB(), id, selectCols...)
}

// FindPilotGP retrieves a single record by ID, and panics on error.
func FindPilotGP(id int, selectCols ...string) *Pilot {
	retobj, err := FindPilot(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindPilot retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPilot(exec boil.Executor, id int, selectCols ...string) (*Pilot, error) {
	pilotObj := &Pilot{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"pilots\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(pilotObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from pilots")
	}

	return pilotObj, nil
}

// FindPilotP retrieves a single record by ID with an executor, and panics on error.
func FindPilotP(exec boil.Executor, id int, selectCols ...string) *Pilot {
	retobj, err := FindPilot(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Pilot) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Pilot) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Pilot) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Pilot) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no pilots provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(pilotColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	pilotInsertCacheMut.RLock()
	cache, cached := pilotInsertCache[key]
	pilotInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			pilotColumns,
			pilotColumnsWithDefault,
			pilotColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(pilotType, pilotMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(pilotType, pilotMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"pilots\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

		if len(cache.retMapping) != 0 {
			cache.query += fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into pilots")
	}

	if !cached {
		pilotInsertCacheMut.Lock()
		pilotInsertCache[key] = cache
		pilotInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Pilot record. See Update for
// whitelist behavior description.
func (o *Pilot) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Pilot record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Pilot) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Pilot, and panics on error.
// See Update for whitelist behavior description.
func (o *Pilot) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Pilot.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Pilot) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	pilotUpdateCacheMut.RLock()
	cache, cached := pilotUpdateCache[key]
	pilotUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(pilotColumns, pilotPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update pilots, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"pilots\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, pilotPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(pilotType, pilotMapping, append(wl, pilotPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update pilots row")
	}

	if !cached {
		pilotUpdateCacheMut.Lock()
		pilotUpdateCache[key] = cache
		pilotUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q pilotQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q pilotQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for pilots")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o PilotSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o PilotSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o PilotSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PilotSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pilotPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"pilots\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(pilotPrimaryKeyColumns), len(colNames)+1, len(pilotPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in pilot slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Pilot) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Pilot) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Pilot) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Pilot) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no pilots provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(pilotColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	pilotUpsertCacheMut.RLock()
	cache, cached := pilotUpsertCache[key]
	pilotUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			pilotColumns,
			pilotColumnsWithDefault,
			pilotColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			pilotColumns,
			pilotPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert pilots, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(pilotPrimaryKeyColumns))
			copy(conflict, pilotPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"pilots\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(pilotType, pilotMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(pilotType, pilotMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert pilots")
	}

	if !cached {
		pilotUpsertCacheMut.Lock()
		pilotUpsertCache[key] = cache
		pilotUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Pilot record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Pilot) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Pilot record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Pilot) DeleteG() error {
	if o == nil {
		return errors.New("models: no Pilot provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Pilot record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Pilot) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Pilot record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Pilot) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Pilot provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), pilotPrimaryKeyMapping)
	sql := "DELETE FROM \"pilots\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from pilots")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q pilotQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q pilotQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no pilotQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from pilots")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o PilotSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o PilotSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Pilot slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o PilotSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PilotSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Pilot slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(pilotBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pilotPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"pilots\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, pilotPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(pilotPrimaryKeyColumns), 1, len(pilotPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from pilot slice")
	}

	if len(pilotAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Pilot) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Pilot) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Pilot) ReloadG() error {
	if o == nil {
		return errors.New("models: no Pilot provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Pilot) Reload(exec boil.Executor) error {
	ret, err := FindPilot(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *PilotSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *PilotSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PilotSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty PilotSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PilotSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	pilots := PilotSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), pilotPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"pilots\".* FROM \"pilots\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, pilotPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(pilotPrimaryKeyColumns), 1, len(pilotPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&pilots)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PilotSlice")
	}

	*o = pilots

	return nil
}

// PilotExists checks if the Pilot row exists.
func PilotExists(exec boil.Executor, id int) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"pilots\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if pilots exists")
	}

	return exists, nil
}

// PilotExistsG checks if the Pilot row exists.
func PilotExistsG(id int) (bool, error) {
	return PilotExists(boil.GetDB(), id)
}

// PilotExistsGP checks if the Pilot row exists. Panics on error.
func PilotExistsGP(id int) bool {
	e, err := PilotExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// PilotExistsP checks if the Pilot row exists. Panics on error.
func PilotExistsP(exec boil.Executor, id int) bool {
	e, err := PilotExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
