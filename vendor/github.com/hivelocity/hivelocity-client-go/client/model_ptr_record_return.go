/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type PtrRecordReturn struct {
	DomainId int32  `json:"domainId,omitempty"`
	Id       int32  `json:"id,omitempty"`
	Address  string `json:"address,omitempty"`
	Type_    string `json:"type,omitempty"`
	Name     string `json:"name,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
	Ttl      int32  `json:"ttl,omitempty"`
}
