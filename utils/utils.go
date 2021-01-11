package utils

//SetupUtils sets up dependant connections for our application
func SetupUtils() error {
	err := setupConfig()
	if err != nil {
		return err
	}

	return nil
}
