/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type CreateCredit struct {
	BillingInfoId int32   `json:"billingInfoId,omitempty"`
	Amount        float32 `json:"amount,omitempty"`
}
