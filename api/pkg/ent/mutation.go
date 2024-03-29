// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/chanzuckerberg/happy/api/pkg/ent/appconfig"
	"github.com/chanzuckerberg/happy/api/pkg/ent/predicate"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeAppConfig = "AppConfig"
)

// AppConfigMutation represents an operation that mutates the AppConfig nodes in the graph.
type AppConfigMutation struct {
	config
	op            Op
	typ           string
	id            *uint
	created_at    *time.Time
	updated_at    *time.Time
	deleted_at    *time.Time
	app_name      *string
	environment   *string
	stack         *string
	key           *string
	value         *string
	source        *appconfig.Source
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*AppConfig, error)
	predicates    []predicate.AppConfig
}

var _ ent.Mutation = (*AppConfigMutation)(nil)

// appconfigOption allows management of the mutation configuration using functional options.
type appconfigOption func(*AppConfigMutation)

// newAppConfigMutation creates new mutation for the AppConfig entity.
func newAppConfigMutation(c config, op Op, opts ...appconfigOption) *AppConfigMutation {
	m := &AppConfigMutation{
		config:        c,
		op:            op,
		typ:           TypeAppConfig,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withAppConfigID sets the ID field of the mutation.
func withAppConfigID(id uint) appconfigOption {
	return func(m *AppConfigMutation) {
		var (
			err   error
			once  sync.Once
			value *AppConfig
		)
		m.oldValue = func(ctx context.Context) (*AppConfig, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().AppConfig.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withAppConfig sets the old AppConfig of the mutation.
func withAppConfig(node *AppConfig) appconfigOption {
	return func(m *AppConfigMutation) {
		m.oldValue = func(context.Context) (*AppConfig, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m AppConfigMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m AppConfigMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of AppConfig entities.
func (m *AppConfigMutation) SetID(id uint) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *AppConfigMutation) ID() (id uint, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *AppConfigMutation) IDs(ctx context.Context) ([]uint, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uint{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().AppConfig.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCreatedAt sets the "created_at" field.
func (m *AppConfigMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *AppConfigMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the AppConfig entity.
// If the AppConfig object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppConfigMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *AppConfigMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the "updated_at" field.
func (m *AppConfigMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the value of the "updated_at" field in the mutation.
func (m *AppConfigMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedAt returns the old "updated_at" field's value of the AppConfig entity.
// If the AppConfig object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppConfigMutation) OldUpdatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdatedAt: %w", err)
	}
	return oldValue.UpdatedAt, nil
}

// ResetUpdatedAt resets all changes to the "updated_at" field.
func (m *AppConfigMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetDeletedAt sets the "deleted_at" field.
func (m *AppConfigMutation) SetDeletedAt(t time.Time) {
	m.deleted_at = &t
}

// DeletedAt returns the value of the "deleted_at" field in the mutation.
func (m *AppConfigMutation) DeletedAt() (r time.Time, exists bool) {
	v := m.deleted_at
	if v == nil {
		return
	}
	return *v, true
}

// OldDeletedAt returns the old "deleted_at" field's value of the AppConfig entity.
// If the AppConfig object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppConfigMutation) OldDeletedAt(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDeletedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDeletedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDeletedAt: %w", err)
	}
	return oldValue.DeletedAt, nil
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (m *AppConfigMutation) ClearDeletedAt() {
	m.deleted_at = nil
	m.clearedFields[appconfig.FieldDeletedAt] = struct{}{}
}

// DeletedAtCleared returns if the "deleted_at" field was cleared in this mutation.
func (m *AppConfigMutation) DeletedAtCleared() bool {
	_, ok := m.clearedFields[appconfig.FieldDeletedAt]
	return ok
}

// ResetDeletedAt resets all changes to the "deleted_at" field.
func (m *AppConfigMutation) ResetDeletedAt() {
	m.deleted_at = nil
	delete(m.clearedFields, appconfig.FieldDeletedAt)
}

// SetAppName sets the "app_name" field.
func (m *AppConfigMutation) SetAppName(s string) {
	m.app_name = &s
}

// AppName returns the value of the "app_name" field in the mutation.
func (m *AppConfigMutation) AppName() (r string, exists bool) {
	v := m.app_name
	if v == nil {
		return
	}
	return *v, true
}

// OldAppName returns the old "app_name" field's value of the AppConfig entity.
// If the AppConfig object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppConfigMutation) OldAppName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldAppName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldAppName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldAppName: %w", err)
	}
	return oldValue.AppName, nil
}

// ResetAppName resets all changes to the "app_name" field.
func (m *AppConfigMutation) ResetAppName() {
	m.app_name = nil
}

// SetEnvironment sets the "environment" field.
func (m *AppConfigMutation) SetEnvironment(s string) {
	m.environment = &s
}

// Environment returns the value of the "environment" field in the mutation.
func (m *AppConfigMutation) Environment() (r string, exists bool) {
	v := m.environment
	if v == nil {
		return
	}
	return *v, true
}

// OldEnvironment returns the old "environment" field's value of the AppConfig entity.
// If the AppConfig object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppConfigMutation) OldEnvironment(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldEnvironment is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldEnvironment requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldEnvironment: %w", err)
	}
	return oldValue.Environment, nil
}

// ResetEnvironment resets all changes to the "environment" field.
func (m *AppConfigMutation) ResetEnvironment() {
	m.environment = nil
}

// SetStack sets the "stack" field.
func (m *AppConfigMutation) SetStack(s string) {
	m.stack = &s
}

// Stack returns the value of the "stack" field in the mutation.
func (m *AppConfigMutation) Stack() (r string, exists bool) {
	v := m.stack
	if v == nil {
		return
	}
	return *v, true
}

// OldStack returns the old "stack" field's value of the AppConfig entity.
// If the AppConfig object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppConfigMutation) OldStack(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStack is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStack requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStack: %w", err)
	}
	return oldValue.Stack, nil
}

// ResetStack resets all changes to the "stack" field.
func (m *AppConfigMutation) ResetStack() {
	m.stack = nil
}

// SetKey sets the "key" field.
func (m *AppConfigMutation) SetKey(s string) {
	m.key = &s
}

// Key returns the value of the "key" field in the mutation.
func (m *AppConfigMutation) Key() (r string, exists bool) {
	v := m.key
	if v == nil {
		return
	}
	return *v, true
}

// OldKey returns the old "key" field's value of the AppConfig entity.
// If the AppConfig object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppConfigMutation) OldKey(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldKey is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldKey requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldKey: %w", err)
	}
	return oldValue.Key, nil
}

// ResetKey resets all changes to the "key" field.
func (m *AppConfigMutation) ResetKey() {
	m.key = nil
}

// SetValue sets the "value" field.
func (m *AppConfigMutation) SetValue(s string) {
	m.value = &s
}

// Value returns the value of the "value" field in the mutation.
func (m *AppConfigMutation) Value() (r string, exists bool) {
	v := m.value
	if v == nil {
		return
	}
	return *v, true
}

// OldValue returns the old "value" field's value of the AppConfig entity.
// If the AppConfig object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppConfigMutation) OldValue(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldValue is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldValue requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldValue: %w", err)
	}
	return oldValue.Value, nil
}

// ResetValue resets all changes to the "value" field.
func (m *AppConfigMutation) ResetValue() {
	m.value = nil
}

// SetSource sets the "source" field.
func (m *AppConfigMutation) SetSource(a appconfig.Source) {
	m.source = &a
}

// Source returns the value of the "source" field in the mutation.
func (m *AppConfigMutation) Source() (r appconfig.Source, exists bool) {
	v := m.source
	if v == nil {
		return
	}
	return *v, true
}

// OldSource returns the old "source" field's value of the AppConfig entity.
// If the AppConfig object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppConfigMutation) OldSource(ctx context.Context) (v appconfig.Source, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSource is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSource requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSource: %w", err)
	}
	return oldValue.Source, nil
}

// ResetSource resets all changes to the "source" field.
func (m *AppConfigMutation) ResetSource() {
	m.source = nil
}

// Where appends a list predicates to the AppConfigMutation builder.
func (m *AppConfigMutation) Where(ps ...predicate.AppConfig) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the AppConfigMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *AppConfigMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.AppConfig, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *AppConfigMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *AppConfigMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (AppConfig).
func (m *AppConfigMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *AppConfigMutation) Fields() []string {
	fields := make([]string, 0, 9)
	if m.created_at != nil {
		fields = append(fields, appconfig.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, appconfig.FieldUpdatedAt)
	}
	if m.deleted_at != nil {
		fields = append(fields, appconfig.FieldDeletedAt)
	}
	if m.app_name != nil {
		fields = append(fields, appconfig.FieldAppName)
	}
	if m.environment != nil {
		fields = append(fields, appconfig.FieldEnvironment)
	}
	if m.stack != nil {
		fields = append(fields, appconfig.FieldStack)
	}
	if m.key != nil {
		fields = append(fields, appconfig.FieldKey)
	}
	if m.value != nil {
		fields = append(fields, appconfig.FieldValue)
	}
	if m.source != nil {
		fields = append(fields, appconfig.FieldSource)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *AppConfigMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case appconfig.FieldCreatedAt:
		return m.CreatedAt()
	case appconfig.FieldUpdatedAt:
		return m.UpdatedAt()
	case appconfig.FieldDeletedAt:
		return m.DeletedAt()
	case appconfig.FieldAppName:
		return m.AppName()
	case appconfig.FieldEnvironment:
		return m.Environment()
	case appconfig.FieldStack:
		return m.Stack()
	case appconfig.FieldKey:
		return m.Key()
	case appconfig.FieldValue:
		return m.Value()
	case appconfig.FieldSource:
		return m.Source()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *AppConfigMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case appconfig.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case appconfig.FieldUpdatedAt:
		return m.OldUpdatedAt(ctx)
	case appconfig.FieldDeletedAt:
		return m.OldDeletedAt(ctx)
	case appconfig.FieldAppName:
		return m.OldAppName(ctx)
	case appconfig.FieldEnvironment:
		return m.OldEnvironment(ctx)
	case appconfig.FieldStack:
		return m.OldStack(ctx)
	case appconfig.FieldKey:
		return m.OldKey(ctx)
	case appconfig.FieldValue:
		return m.OldValue(ctx)
	case appconfig.FieldSource:
		return m.OldSource(ctx)
	}
	return nil, fmt.Errorf("unknown AppConfig field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *AppConfigMutation) SetField(name string, value ent.Value) error {
	switch name {
	case appconfig.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case appconfig.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case appconfig.FieldDeletedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDeletedAt(v)
		return nil
	case appconfig.FieldAppName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetAppName(v)
		return nil
	case appconfig.FieldEnvironment:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnvironment(v)
		return nil
	case appconfig.FieldStack:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStack(v)
		return nil
	case appconfig.FieldKey:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetKey(v)
		return nil
	case appconfig.FieldValue:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetValue(v)
		return nil
	case appconfig.FieldSource:
		v, ok := value.(appconfig.Source)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSource(v)
		return nil
	}
	return fmt.Errorf("unknown AppConfig field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *AppConfigMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *AppConfigMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *AppConfigMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown AppConfig numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *AppConfigMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(appconfig.FieldDeletedAt) {
		fields = append(fields, appconfig.FieldDeletedAt)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *AppConfigMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *AppConfigMutation) ClearField(name string) error {
	switch name {
	case appconfig.FieldDeletedAt:
		m.ClearDeletedAt()
		return nil
	}
	return fmt.Errorf("unknown AppConfig nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *AppConfigMutation) ResetField(name string) error {
	switch name {
	case appconfig.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case appconfig.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case appconfig.FieldDeletedAt:
		m.ResetDeletedAt()
		return nil
	case appconfig.FieldAppName:
		m.ResetAppName()
		return nil
	case appconfig.FieldEnvironment:
		m.ResetEnvironment()
		return nil
	case appconfig.FieldStack:
		m.ResetStack()
		return nil
	case appconfig.FieldKey:
		m.ResetKey()
		return nil
	case appconfig.FieldValue:
		m.ResetValue()
		return nil
	case appconfig.FieldSource:
		m.ResetSource()
		return nil
	}
	return fmt.Errorf("unknown AppConfig field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *AppConfigMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *AppConfigMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *AppConfigMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *AppConfigMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *AppConfigMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *AppConfigMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *AppConfigMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown AppConfig unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *AppConfigMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown AppConfig edge %s", name)
}
