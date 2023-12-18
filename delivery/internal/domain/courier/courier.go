package courier

type Courier struct {
	uuid string
	name string
}

func NewCourier(
	uuid string,
	name string,
) (*Courier, error) {
	return &Courier{
		uuid: uuid,
		name: name,
	}, nil
}

func (c Courier) UUID() string {
	return c.uuid
}

func (c Courier) Name() string {
	return c.name
}

func UnmarshalCourierFromDatabase(
	uuid string,
	name string,
) (*Courier, error) {

	return &Courier{
		uuid: uuid,
		name: name,
	}, nil
}
