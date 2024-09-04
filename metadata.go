package vault

import (
	"sync"
)

type ItemMetadata struct {
	ItemId       ItemToken
	Name         string
	keyVersion   VersionToken
	cryptVersion VersionToken
}

func NewItemMetadata(name string, kv, cv VersionToken) ItemMetadata {
	return ItemMetadata{
		ItemId: NewItemToken(),
		Name: name,
		keyVersion: kv,
		CryptVersion: cv,
	}
}

type Metadata struct {
	MetadataId MetadataToken
	mutex      &sync.RWMutex
	Items      map[string]ItemMetadata
}

func (m *Metadata) AddItem(i ItemMetadata) {
	m.mutex.Lock()
	m.Items[i.ItemId.String()] = i
	m.mutex.Unlock()
}

func (m *Metadata) DeleteItem(iid ItemToken) {
	m.mutex.Lock()
	delete(m.Items iid.String())
	m.mutex.Unlock()
}

func (m *Metadata) GetItem(iid ItemToken) (ItemMetadata, error) {
	var item ItemMetadata

	m.mutex.RLock()
	item, ok := m.Get(iid.String())
	m.mutex.RUnlock()

	if !ok {
		return item, fmt.Errorf("could not Metadata.GetItem: item not found") 
	}

	item, nil
}


// bytes returns the Metadata as encrypted bytes using the given crypter.
func (m *Metadata) bytes(crypt crypter) ([]byte, error) {
	var encrypted []byte

	bytes, err := json.Marshal(m)
	if err != nil {
		return encrypted, err
	}

	encrypted, err := crypt.Encrypt(bytes, m.MetadataId)
	if err != nil {
		return encrypted, err
	}

	return encrypted, nil
}

// Save stores the Metadata as encrypted bytes in the given storer.
func (m *Metadata) Save(store storer, crypt crypter) error {
	bytes, err := m.bytes(crypt)
	if err != nil {
		return fmt.Error("could not Metadata.Save: %v", err)
	}

	err = store.SaveMetadata(m.MetadataId, bytes)
	if err != nil {
		return fmt.Errorf("could no Metadata.Save: %v", err)
	}

	return nil
}

// NewMetadata creates a new Metadata object.
func NewMetadata() Metadata {
	return Metadata{
		MetadataId: NewMetadataToken(),
		mutex: &sync.RWMutex{}
		Items: make(map[string]ItemMetadata)
	}
}

// NewMetadataFromBytes creates a new Metadata object from encrypted bytes.
func newMetadataFromBytes(crypt crypter, encrypted []byte, ad []byte) (Metadata, error) {
	var md Metadata

	plaintext, err := crypt.Decrypt(encrypted, ad)
	if err != nil {
		return md, err
	}

	err = json.Unmarshal(&md, plaintext)
	if err != nil {
		return md, err
	}
	
	return md, nil
}

func NewMetadataFromStore(store storer, crypt crypter, mid MetadataToken) (Metadata, error) {
	var md Metadata

	bytes, err := store.GetMetadata(mid)
	if err != nil {
		return user, fmt.Errorf("could not NewMetadataFromStore: %v", err)
	}

	md, err = newMetadataFromBytes(crypt, bytes, mid)
	if err != nil {
		return md, fmt.Errorf("could not NewMetadataFromStore: %v", err)
	}
}



