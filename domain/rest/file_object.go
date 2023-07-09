package rest_object

import rest "github.com/brutalzinn/go-easy-rest"

type FileObject struct {
	Name string
	Link rest.Link
}

func (restObject RestObject) CreateFileObject() rest.Resource {
	resourceData := rest.Resource{
		Object:    FILE_OBJECT,
		Attribute: restObject.Attribute,
		Link:      restObject.Link,
	}
	return resourceData
}
