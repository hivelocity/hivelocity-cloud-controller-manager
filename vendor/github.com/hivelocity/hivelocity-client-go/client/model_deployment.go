/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Deployment struct {
	Empty                   bool          `json:"empty,omitempty"`
	DeploymentName          string        `json:"deploymentName,omitempty"`
	StartedProvisioning     bool          `json:"startedProvisioning,omitempty"`
	DeploymentId            int32         `json:"deploymentId,omitempty"`
	OrderNumber             string        `json:"orderNumber,omitempty"`
	Price                   float32       `json:"price,omitempty"`
	DeploymentConfiguration []interface{} `json:"deploymentConfiguration,omitempty"`
}
