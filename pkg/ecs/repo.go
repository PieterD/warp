package ecs

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"time"
)

type Repository struct {
	rand       *rand.Rand
	entities   EntityIDSet
	components map[reflect.Type]map[ID]reflect.Value
}

func NewRepository() *Repository {
	return &Repository{
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
		entities:   NewEntityIDSet(),
		components: make(map[reflect.Type]map[ID]reflect.Value),
	}
}

func (repo *Repository) NewEntity() Entity {
	for n := 0; n < 256; n++ {
		randomId := ID(repo.rand.Int63())
		if _, ok := repo.entities[randomId]; ok {
			continue
		}
		repo.entities[randomId] = struct{}{}
		return newEntity(repo, randomId)
	}
	panic(fmt.Errorf("took too long to find unused id"))
}

// Entities visits all entities that possess each of the provided components
// in ID order, and calls f for each Entity in the Repository.
// Any of the rawComponents may be a straight component or a pointer to one.
// Any rawComponent pointer will be set to the value of the Entity passed to f.
// Unless it is nil, in which case it will only allow instances of Entity without that component.
func (repo *Repository) Entities(f func(Entity) error, rawComponents ...interface{}) error {
	type entityFilter struct {
		inverted      bool
		set           bool
		entity        bool
		typ           reflect.Type
		val           reflect.Value
		componentByID map[ID]reflect.Value
		size          int
	}
	var filters []entityFilter
	for _, raw := range rawComponents {
		val := reflect.ValueOf(raw)
		typ := val.Type()
		kind := typ.Kind()
		if kind == reflect.Ptr {
			if val.IsNil() {
				filters = append(filters, entityFilter{
					inverted: true,
					typ:      typ.Elem(),
				})
				continue
			}
			filters = append(filters, entityFilter{
				set: true,
				typ: typ.Elem(),
				val: val.Elem(),
			})
			continue
		}
		filters = append(filters, entityFilter{
			typ: typ,
		})
	}
	const maxInt = int((^uint(0)) >> 1)
	for i := range filters {
		filter := &filters[i]
		componentByID, ok := repo.components[filter.typ]
		if !ok {
			if filter.inverted {
				filter.size = maxInt
				continue
			}
			//return fmt.Errorf("repository does not contain filtering entity type: %v", filter.typ)
			return nil
		}
		filter.componentByID = componentByID
		filter.size = len(componentByID)
		if filter.size == 0 {
			if filter.inverted {
				filter.size = maxInt
				continue
			}
			return nil
		}
	}
	filters = append(filters, entityFilter{
		entity: true,
		size:   len(repo.entities),
	})
	sort.SliceStable(filters, func(i, j int) bool {
		return filters[i].size < filters[j].size
	})
	if filters[0].inverted {
		return fmt.Errorf("no valid non-inverted filters provided")
	}
	process := func(id ID) error {
		components := make([]reflect.Value, 0, len(filters))
		for _, filter := range filters {
			componentRaw, ok := filter.componentByID[id]
			if !ok {
				if !filter.inverted {
					break
				}
				components = append(components, reflect.Value{})
				continue
			}
			if filter.inverted {
				break
			}
			components = append(components, componentRaw)
		}
		if len(components) != len(filters) {
			return nil
		}
		for i := range components {
			component := components[i]
			filter := filters[i]
			if filter.set {
				filter.val.Set(component)
			}
		}
		if err := f(newEntity(repo, id)); err != nil {
			return err
		}
		return nil
	}
	if filters[0].entity {
		for id := range repo.entities {
			if err := process(id); err != nil {
				return fmt.Errorf("processing %v: %w", id, err)
			}
		}
	}
	for id := range filters[0].componentByID {
		if err := process(id); err != nil {
			return fmt.Errorf("processing %v: %w", id, err)
		}
	}
	return nil
}

func (repo *Repository) HasEntity(id ID) bool {
	_, ok := repo.entities[id]
	return ok
}

func (repo *Repository) DelEntity(id ID) {
	delete(repo.entities, id)
	for _, components := range repo.components {
		delete(components, id)
	}
}
