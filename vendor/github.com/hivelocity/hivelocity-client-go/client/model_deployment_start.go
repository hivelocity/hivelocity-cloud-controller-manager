/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DeploymentStart struct {
	BillingInfo int32  `json:"billingInfo"`
	Script      string `json:"script,omitempty"`
}
