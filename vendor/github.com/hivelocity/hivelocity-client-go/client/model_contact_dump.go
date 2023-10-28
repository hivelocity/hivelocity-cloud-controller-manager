/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ContactDump struct {
	Description string `json:"description,omitempty"`
	ClientId    int32  `json:"clientId,omitempty"`
	Email       string `json:"email"`
	ContactId   int32  `json:"contactId,omitempty"`
	IsClient    bool   `json:"isClient,omitempty"`
	Active      int32  `json:"active"`
	FullName    string `json:"fullName"`
	Phone       string `json:"phone,omitempty"`
}
