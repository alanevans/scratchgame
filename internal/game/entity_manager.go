package game

import (
	"fmt"
	"sync"
)

// EntityID represents a unique identifier for entities
type EntityID uint32

// Component interface that all components must implement
type Component interface {
	GetType() ComponentType
}

// ComponentType represents different component types
type ComponentType int

const (
	ComponentTypePosition ComponentType = iota
	ComponentTypeVelocity
	ComponentTypeSprite
	ComponentTypeCollider
	ComponentTypePlayer
	ComponentTypePlatform
	ComponentTypeGravity
	ComponentTypeInput
)

// EntityManager manages all entities and their components
type EntityManager struct {
	nextEntityID EntityID
	entities     map[EntityID]bool
	components   map[ComponentType]map[EntityID]Component
	mutex        sync.RWMutex
}

// NewEntityManager creates a new entity manager
func NewEntityManager() *EntityManager {
	return &EntityManager{
		nextEntityID: 1,
		entities:     make(map[EntityID]bool),
		components:   make(map[ComponentType]map[EntityID]Component),
	}
}

// CreateEntity creates a new entity and returns its ID
func (em *EntityManager) CreateEntity() EntityID {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	id := em.nextEntityID
	em.nextEntityID++
	em.entities[id] = true

	return id
}

// DestroyEntity removes an entity and all its components
func (em *EntityManager) DestroyEntity(entityID EntityID) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Remove from entities map
	delete(em.entities, entityID)

	// Remove all components for this entity
	for componentType := range em.components {
		delete(em.components[componentType], entityID)
	}
}

// AddComponent adds a component to an entity
func (em *EntityManager) AddComponent(entityID EntityID, component Component) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	componentType := component.GetType()

	if em.components[componentType] == nil {
		em.components[componentType] = make(map[EntityID]Component)
	}

	em.components[componentType][entityID] = component
}

// GetComponent retrieves a component from an entity
func (em *EntityManager) GetComponent(entityID EntityID, componentType ComponentType) (Component, bool) {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	if components, exists := em.components[componentType]; exists {
		if component, exists := components[entityID]; exists {
			return component, true
		}
	}

	return nil, false
}

// HasComponent checks if an entity has a specific component
func (em *EntityManager) HasComponent(entityID EntityID, componentType ComponentType) bool {
	_, exists := em.GetComponent(entityID, componentType)
	return exists
}

// GetEntitiesWithComponent returns all entities that have a specific component
func (em *EntityManager) GetEntitiesWithComponent(componentType ComponentType) []EntityID {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	var entities []EntityID

	if components, exists := em.components[componentType]; exists {
		for entityID := range components {
			if em.entities[entityID] {
				entities = append(entities, entityID)
			}
		}
	}

	return entities
}

// GetEntitiesWithComponents returns all entities that have all specified components
func (em *EntityManager) GetEntitiesWithComponents(componentTypes ...ComponentType) []EntityID {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	var entities []EntityID

	for entityID := range em.entities {
		hasAll := true
		for _, componentType := range componentTypes {
			if !em.hasComponentUnsafe(entityID, componentType) {
				hasAll = false
				break
			}
		}
		if hasAll {
			entities = append(entities, entityID)
		}
	}

	return entities
}

// hasComponentUnsafe checks if an entity has a component without locking (internal use)
func (em *EntityManager) hasComponentUnsafe(entityID EntityID, componentType ComponentType) bool {
	if components, exists := em.components[componentType]; exists {
		_, exists := components[entityID]
		return exists
	}
	return false
}

// PrintDebugInfo prints debug information about entities and components
func (em *EntityManager) PrintDebugInfo() {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	fmt.Printf("Total entities: %d\n", len(em.entities))
	for componentType, components := range em.components {
		fmt.Printf("Component type %d: %d entities\n", componentType, len(components))
	}
}
