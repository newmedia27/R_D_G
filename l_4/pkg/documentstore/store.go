package documentstore

type Store struct {
	Collections map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		Collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	// Створюємо нову колекцію і повертаємо `true` якщо колекція була створена
	// Якщо ж колекція вже створеня то повертаємо `false` та nil
	if name == "" {
		return false, nil
	}
	if _, ok := s.Collections[name]; ok {
		return false, nil
	}

	col := &Collection{
		Documents: make(map[string]*Document),
		name:      name,
		cfg:       *cfg,
	}
	s.Collections[name] = col

	return true, col
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	if col, ok := s.Collections[name]; ok {
		return col, true
	}
	return nil, false
}

func (s *Store) DeleteCollection(name string) bool {
	if _, ok := s.Collections[name]; ok {
		delete(s.Collections, name)
		return true
	}
	return false
}
