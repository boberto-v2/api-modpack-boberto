package rest_object

import (
	rest "github.com/brutalzinn/go-easy-rest"
)

type UserCredentialsObject struct {
	AccessToken string `json:"access_token"`
}

func (restObject RestObject) CreateUserCredentialsObject(data string) RestObject {
	restObject.Resource = rest.Resource{
		Object: USER_CREDENTIAL_OBJECT,
		Attribute: UserCredentialsObject{
			AccessToken: data,
		},
		Link: restObject.Link,
	}
	restObject.create()
	return restObject
}
