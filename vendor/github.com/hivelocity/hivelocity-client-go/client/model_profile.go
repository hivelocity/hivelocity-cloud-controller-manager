/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Profile struct {
	State    interface{} `json:"state,omitempty"`
	City     interface{} `json:"city,omitempty"`
	FullName interface{} `json:"fullName,omitempty"`
	IsClient bool        `json:"isClient,omitempty"`
	Email    string      `json:"email,omitempty"`
	Created  interface{} `json:"created,omitempty"`
	Company  interface{} `json:"company,omitempty"`
	Zip      interface{} `json:"zip,omitempty"`
	Last     string      `json:"last,omitempty"`
	Id       int32       `json:"id,omitempty"`
	Fax      interface{} `json:"fax,omitempty"`
	Address  interface{} `json:"address,omitempty"`
	MetaData interface{} `json:"metaData,omitempty"`
	First    string      `json:"first,omitempty"`
	Phone    string      `json:"phone,omitempty"`
	Login    string      `json:"login,omitempty"`
	Country  interface{} `json:"country,omitempty"`
}