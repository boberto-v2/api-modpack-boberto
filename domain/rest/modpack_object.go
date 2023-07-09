package rest_object

import (
	rest "github.com/brutalzinn/go-easy-rest"
)

func (restObject RestObject) CreateModPackObject() rest.Resource {
	resource := rest.Resource{
		Object:    MODPACK_OBJECT,
		Attribute: restObject.Attribute,
		Link:      restObject.Link,
	}
	return resource
}
