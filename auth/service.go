package auth

// Service defines auth behaviors
type Service interface {
	// RegisterByPassword gets a new user by username and password
	RegisterByPassword(username, pass string) (int, error)

	// RegisterByMobileCode gets a new user by mobile and code
	RegisterByMobileCode(mobile, code string) (int, error)

	// RegisterByEmailCode gets a new user by email and code
	RegisterByEmailCode(email, code string) (int, error)

	// LoginByPassword login system by username and password
	// and return user id if authentication success
	LoginByPassword(username, pass string) (int, error)

	// LoginByMobileCode login system by mobile number and validate code
	// and return user id if authentication success
	LoginByMobileCode(mobile, code string) (int, error)

	// LoginByEmailCode login system by email address and validate code
	// and return user id if authentication success
	LoginByEmailCode(email, code string) (int, error)
}

// CodeValidator defences behaviors of a code validator
type CodeValidator interface {
	// Validate validates whether the code matched with the key
	Validate(key, code string) error
}

// service is an internal implement of Service interface
type service struct {
	repo Repo
	v    CodeValidator
}

func (s service) RegisterByPassword(username, pass string) (int, error) {
	return s.repo.RegisterByPassword(username, pass)
}

func (s service) RegisterByMobileCode(mobile, code string) (int, error) {
	if err := s.v.Validate(mobile, code); err != nil {
		return 0, err
	}

	return s.repo.RegisterByMobile(mobile)
}

func (s service) RegisterByEmailCode(email, code string) (int, error) {
	if err := s.v.Validate(email, code); err != nil {
		return 0, err
	}

	return s.repo.RegisterByEmail(email)
}

func (s service) LoginByPassword(username, pass string) (int, error) {
	return s.repo.LoginByPassword(username, pass)
}

func (s service) LoginByMobileCode(mobile, code string) (int, error) {
	if err := s.v.Validate(mobile, code); err != nil {
		return 0, err
	}

	return s.repo.LoginByMobile(mobile)
}

func (s service) LoginByEmailCode(email, code string) (int, error) {
	if err := s.v.Validate(email, code); err != nil {
		return 0, err
	}

	return s.repo.LoginByEmail(email)
}
