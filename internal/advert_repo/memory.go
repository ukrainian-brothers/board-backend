package advert_repo

import (
	"github.com/google/uuid"
	"github.com/ukrainian-brothers/board-backend/domain/advert"
	"sync"
)

type MemoryAdvertRepository struct {
	memory map[uuid.UUID]*advert.Advert
	sync.Mutex
}

func NewMemoryAdvertRepository(memory map[uuid.UUID]*advert.Advert) *MemoryAdvertRepository {
	return &MemoryAdvertRepository{memory: memory}
}


func (repo *MemoryAdvertRepository) Get(id uuid.UUID) (*advert.Advert, error) {
	if adv, ok := repo.memory[id]; ok {
		return adv, nil
	}
	return nil, AdvertNotFound
}

func (repo *MemoryAdvertRepository) Add(advert *advert.Advert) error {
	if _, ok := repo.memory[advert.ID]; ok {
		return AdvertAlreadyExists
	}
	repo.Lock()
	repo.memory[advert.ID] = advert
	repo.Unlock()
	return nil
}

func (repo *MemoryAdvertRepository) Delete(id uuid.UUID) error {
	if _, ok := repo.memory[id]; !ok {
		return AdvertNotFound
	}
	repo.Lock()
	delete(repo.memory, id)
	repo.Unlock()
	return nil
}