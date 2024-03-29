/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Location struct {
	// true|false if edge site.
	Edge bool `json:"edge,omitempty"`
	// The unique facility code.
	Code string `json:"code,omitempty"`
	// The unique facility name.
	Title    string            `json:"title,omitempty"`
	Location *LocationLocation `json:"location,omitempty"`
	// true|false if core site.
	Core bool `json:"core,omitempty"`
}
