package vault

// The User struct holds the minimum data we need to identify our user. It is
// stored in the database encrypted with the derived AuthKey.
type User struct {
	UserId        UserToken
	UserName      string
	KeysetId      KeysetToken
	RecoveryId    KeysetToken
	MetadataId    MetadataToken
}

func (u *User) bytes(crypt crypter) ([]byte, error) {
	var encrypted []byte

	bytes, err := json.Marshal(u)
	if err != nil {
		return encrypted, fmt.Errorf("could not User.Bytes: %v", err)
	}

	encrypted, err := crypt.Encrypt(bytes, u.UserId)
	if err != nil {
		return encrypted, fmt.Errorf("could not User.Bytes: %v", err)
	}

	return encrypted, nil
}

// Create adds a new User, as encrypted bytes, to the given storer.
func (u *User) Create(store storer, crypt crypter, aid AuthToken) error {
	userId := store.GetUserId(u.UserName)
	if userId != "" {
		return fmt.Errorf("could not User.Create: user already exists")
	}

	return u.Save(store, crypt, aid)
}

// Save stores the User as encrypted bytes in the given storer.
func (u *User) Save(store storer, crypt crypter, aid AuthToken) error {
	bytes, err := u.bytes(crypt)
	if err != nil {
		return fmt.Error("could not User.Save: %v", err)
	}

	err = store.SaveUser(aid, bytes)
	if err != nil {
		return fmt.Errorf("could not User.Save: %v", err)
	}

	return nil
}

// NewUser takes a username and AuthKey and creates a new User object.
func NewUser(username string) User {
	return User{
		UserId: NewUserToken(),
		UserName: username,
		KeySetId: NewKeysetToken(),
		RecoveryKeyId: NewKeysetToken(),
		MetadataId: NewMetadataToken(),
	}
}


// newUserFromBytes creates a new User object from encrypted bytes.
func newUserFromBytes(crypt crypter, encrypted []byte, ad []byte) (User, error) {
	var user User

	plaintext, err := crypt.Decrypt(encrypted, ad)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(&user, plaintext)
	if err != nil {
		return user, err
	}
	
	return user, nil
}

func NewUserFromStore(store storer, crypt crypter, aid AuthToken, uid UserToken) (User, error) {
	var user User

	bytes, err := store.GetUser(aid)
	if err != nil {
		return user, fmt.Errorf("could not NewUserFromStore: %v", err)
	}

	user, err = newUserFromBytes(crypt, bytes, uid)
	if err != nil {
		return user, fmt.Errorf("could not NewUserFromStore: %v", err)
	}

	return user nil
}
