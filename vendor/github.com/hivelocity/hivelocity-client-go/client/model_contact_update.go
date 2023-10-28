/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ContactUpdate struct {
	Description string `json:"description,omitempty"`
	Email       string `json:"email,omitempty"`
	Active      int32  `json:"active,omitempty"`
	FullName    string `json:"fullName,omitempty"`
	Phone       string `json:"phone,omitempty"`
}
