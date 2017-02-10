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

// Jet is an object representing the database table.
type Jet struct {
	ID      int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	PilotID int    `boil:"pilot_id" json:"pilot_id" toml:"pilot_id" yaml:"pilot_id"`
	Age     int    `boil:"age" json:"age" toml:"age" yaml:"age"`
	Name    string `boil:"name" json:"name" toml:"name" yaml:"name"`
	Color   string `boil:"color" json:"color" toml:"color" yaml:"color"`

	R *jetR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L jetL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// jetR is where relationships are stored.
type jetR struct {
	Pilot *Pilot
}

// jetL is where Load methods for each relationship are stored.
type jetL struct{}

var (
	jetColumns               = []string{"id", "pilot_id", "age", "name", "color"}
	jetColumnsWithoutDefault = []string{"pilot_id", "age", "name", "color"}
	jetColumnsWithDefault    = []string{"id"}
	jetPrimaryKeyColumns     = []string{"id"}
)

type (
	// JetSlice is an alias for a slice of pointers to Jet.
	// This should generally be used opposed to []Jet.
	JetSlice []*Jet
	// JetHook is the signature for custom Jet hook methods
	JetHook func(boil.Executor, *Jet) error

	jetQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	jetType                 = reflect.TypeOf(&Jet{})
	jetMapping              = queries.MakeStructMapping(jetType)
	jetPrimaryKeyMapping, _ = queries.BindMapping(jetType, jetMapping, jetPrimaryKeyColumns)
	jetInsertCacheMut       sync.RWMutex
	jetInsertCache          = make(map[string]insertCache)
	jetUpdateCacheMut       sync.RWMutex
	jetUpdateCache          = make(map[string]updateCache)
	jetUpsertCacheMut       sync.RWMutex
	jetUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var jetBeforeInsertHooks []JetHook
var jetBeforeUpdateHooks []JetHook
var jetBeforeDeleteHooks []JetHook
var jetBeforeUpsertHooks []JetHook

var jetAfterInsertHooks []JetHook
var jetAfterSelectHooks []JetHook
var jetAfterUpdateHooks []JetHook
var jetAfterDeleteHooks []JetHook
var jetAfterUpsertHooks []JetHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Jet) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range jetBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Jet) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range jetBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Jet) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range jetBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Jet) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range jetBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Jet) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range jetAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Jet) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range jetAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Jet) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range jetAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Jet) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range jetAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Jet) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range jetAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddJetHook registers your hook function for all future operations.
func AddJetHook(hookPoint boil.HookPoint, jetHook JetHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		jetBeforeInsertHooks = append(jetBeforeInsertHooks, jetHook)
	case boil.BeforeUpdateHook:
		jetBeforeUpdateHooks = append(jetBeforeUpdateHooks, jetHook)
	case boil.BeforeDeleteHook:
		jetBeforeDeleteHooks = append(jetBeforeDeleteHooks, jetHook)
	case boil.BeforeUpsertHook:
		jetBeforeUpsertHooks = append(jetBeforeUpsertHooks, jetHook)
	case boil.AfterInsertHook:
		jetAfterInsertHooks = append(jetAfterInsertHooks, jetHook)
	case boil.AfterSelectHook:
		jetAfterSelectHooks = append(jetAfterSelectHooks, jetHook)
	case boil.AfterUpdateHook:
		jetAfterUpdateHooks = append(jetAfterUpdateHooks, jetHook)
	case boil.AfterDeleteHook:
		jetAfterDeleteHooks = append(jetAfterDeleteHooks, jetHook)
	case boil.AfterUpsertHook:
		jetAfterUpsertHooks = append(jetAfterUpsertHooks, jetHook)
	}
}

// OneP returns a single jet record from the query, and panics on error.
func (q jetQuery) OneP() *Jet {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single jet record from the query.
func (q jetQuery) One() (*Jet, error) {
	o := &Jet{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for jets")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Jet records from the query, and panics on error.
func (q jetQuery) AllP() JetSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Jet records from the query.
func (q jetQuery) All() (JetSlice, error) {
	var o JetSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Jet slice")
	}

	if len(jetAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Jet records in the query, and panics on error.
func (q jetQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Jet records in the query.
func (q jetQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count jets rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q jetQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q jetQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if jets exists")
	}

	return count > 0, nil
}

// PilotG pointed to by the foreign key.
func (o *Jet) PilotG(mods ...qm.QueryMod) pilotQuery {
	return o.Pilot(boil.GetDB(), mods...)
}

// Pilot pointed to by the foreign key.
func (o *Jet) Pilot(exec boil.Executor, mods ...qm.QueryMod) pilotQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.PilotID),
	}

	queryMods = append(queryMods, mods...)

	query := Pilots(exec, queryMods...)
	queries.SetFrom(query.Query, "\"pilots\"")

	return query
}

// LoadPilot allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (jetL) LoadPilot(e boil.Executor, singular bool, maybeJet interface{}) error {
	var slice []*Jet
	var object *Jet

	count := 1
	if singular {
		object = maybeJet.(*Jet)
	} else {
		slice = *maybeJet.(*JetSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &jetR{}
		}
		args[0] = object.PilotID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &jetR{}
			}
			args[i] = obj.PilotID
		}
	}

	query := fmt.Sprintf(
		"select * from \"pilots\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Pilot")
	}
	defer results.Close()

	var resultSlice []*Pilot
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Pilot")
	}

	if len(jetAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Pilot = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.PilotID == foreign.ID {
				local.R.Pilot = foreign
				break
			}
		}
	}

	return nil
}

// SetPilotG of the jet to the related item.
// Sets o.R.Pilot to related.
// Adds o to related.R.Jets.
// Uses the global database handle.
func (o *Jet) SetPilotG(insert bool, related *Pilot) error {
	return o.SetPilot(boil.GetDB(), insert, related)
}

// SetPilotP of the jet to the related item.
// Sets o.R.Pilot to related.
// Adds o to related.R.Jets.
// Panics on error.
func (o *Jet) SetPilotP(exec boil.Executor, insert bool, related *Pilot) {
	if err := o.SetPilot(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetPilotGP of the jet to the related item.
// Sets o.R.Pilot to related.
// Adds o to related.R.Jets.
// Uses the global database handle and panics on error.
func (o *Jet) SetPilotGP(insert bool, related *Pilot) {
	if err := o.SetPilot(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetPilot of the jet to the related item.
// Sets o.R.Pilot to related.
// Adds o to related.R.Jets.
func (o *Jet) SetPilot(exec boil.Executor, insert bool, related *Pilot) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"jets\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"pilot_id"}),
		strmangle.WhereClause("\"", "\"", 2, jetPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.PilotID = related.ID

	if o.R == nil {
		o.R = &jetR{
			Pilot: related,
		}
	} else {
		o.R.Pilot = related
	}

	if related.R == nil {
		related.R = &pilotR{
			Jets: JetSlice{o},
		}
	} else {
		related.R.Jets = append(related.R.Jets, o)
	}

	return nil
}

// JetsG retrieves all records.
func JetsG(mods ...qm.QueryMod) jetQuery {
	return Jets(boil.GetDB(), mods...)
}

// Jets retrieves all the records using an executor.
func Jets(exec boil.Executor, mods ...qm.QueryMod) jetQuery {
	mods = append(mods, qm.From("\"jets\""))
	return jetQuery{NewQuery(exec, mods...)}
}

// FindJetG retrieves a single record by ID.
func FindJetG(id int, selectCols ...string) (*Jet, error) {
	return FindJet(boil.GetDB(), id, selectCols...)
}

// FindJetGP retrieves a single record by ID, and panics on error.
func FindJetGP(id int, selectCols ...string) *Jet {
	retobj, err := FindJet(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindJet retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindJet(exec boil.Executor, id int, selectCols ...string) (*Jet, error) {
	jetObj := &Jet{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"jets\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(jetObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from jets")
	}

	return jetObj, nil
}

// FindJetP retrieves a single record by ID with an executor, and panics on error.
func FindJetP(exec boil.Executor, id int, selectCols ...string) *Jet {
	retobj, err := FindJet(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Jet) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Jet) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Jet) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Jet) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no jets provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(jetColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	jetInsertCacheMut.RLock()
	cache, cached := jetInsertCache[key]
	jetInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			jetColumns,
			jetColumnsWithDefault,
			jetColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(jetType, jetMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(jetType, jetMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"jets\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into jets")
	}

	if !cached {
		jetInsertCacheMut.Lock()
		jetInsertCache[key] = cache
		jetInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Jet record. See Update for
// whitelist behavior description.
func (o *Jet) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Jet record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Jet) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Jet, and panics on error.
// See Update for whitelist behavior description.
func (o *Jet) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Jet.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Jet) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	jetUpdateCacheMut.RLock()
	cache, cached := jetUpdateCache[key]
	jetUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(jetColumns, jetPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update jets, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"jets\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, jetPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(jetType, jetMapping, append(wl, jetPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update jets row")
	}

	if !cached {
		jetUpdateCacheMut.Lock()
		jetUpdateCache[key] = cache
		jetUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q jetQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q jetQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for jets")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o JetSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o JetSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o JetSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o JetSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), jetPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"jets\" SET %s WHERE (\"id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(jetPrimaryKeyColumns), len(colNames)+1, len(jetPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in jet slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Jet) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Jet) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Jet) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Jet) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no jets provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(jetColumnsWithDefault, o)

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

	jetUpsertCacheMut.RLock()
	cache, cached := jetUpsertCache[key]
	jetUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			jetColumns,
			jetColumnsWithDefault,
			jetColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			jetColumns,
			jetPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert jets, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(jetPrimaryKeyColumns))
			copy(conflict, jetPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"jets\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(jetType, jetMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(jetType, jetMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert jets")
	}

	if !cached {
		jetUpsertCacheMut.Lock()
		jetUpsertCache[key] = cache
		jetUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Jet record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Jet) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Jet record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Jet) DeleteG() error {
	if o == nil {
		return errors.New("models: no Jet provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Jet record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Jet) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Jet record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Jet) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Jet provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), jetPrimaryKeyMapping)
	sql := "DELETE FROM \"jets\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from jets")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q jetQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q jetQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no jetQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from jets")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o JetSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o JetSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Jet slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o JetSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o JetSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Jet slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(jetBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), jetPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"jets\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, jetPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(jetPrimaryKeyColumns), 1, len(jetPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from jet slice")
	}

	if len(jetAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Jet) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Jet) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Jet) ReloadG() error {
	if o == nil {
		return errors.New("models: no Jet provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Jet) Reload(exec boil.Executor) error {
	ret, err := FindJet(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *JetSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *JetSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *JetSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty JetSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *JetSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	jets := JetSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), jetPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"jets\".* FROM \"jets\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, jetPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(jetPrimaryKeyColumns), 1, len(jetPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&jets)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in JetSlice")
	}

	*o = jets

	return nil
}

// JetExists checks if the Jet row exists.
func JetExists(exec boil.Executor, id int) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"jets\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if jets exists")
	}

	return exists, nil
}

// JetExistsG checks if the Jet row exists.
func JetExistsG(id int) (bool, error) {
	return JetExists(boil.GetDB(), id)
}

// JetExistsGP checks if the Jet row exists. Panics on error.
func JetExistsGP(id int) bool {
	e, err := JetExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// JetExistsP checks if the Jet row exists. Panics on error.
func JetExistsP(exec boil.Executor, id int) bool {
	e, err := JetExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
