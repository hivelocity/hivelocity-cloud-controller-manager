/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type OrderGroupCreate struct {
	SameRack bool   `json:"same_rack,omitempty"`
	Name     string `json:"name"`
}
