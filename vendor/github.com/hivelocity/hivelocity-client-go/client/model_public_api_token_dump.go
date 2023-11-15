/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type PublicApiTokenDump struct {
	Token       string            `json:"token"`
	IpAddresses *PublicApiTokenIp `json:"ipAddresses,omitempty"`
	Name        string            `json:"name,omitempty"`
}