/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DeploymentCustomization struct {
	// Contents of iPXE script if not supplying URL
	CustomIPXEScriptContents string   `json:"customIPXEScriptContents,omitempty"`
	LocationCode             string   `json:"locationCode,omitempty"`
	Hostnames                []string `json:"hostnames"`
	// Specify Ignition file ID for CoreOS/Flatcar provisions
	IgnitionIds []int32 `json:"ignitionIds,omitempty"`
	// URL to download custom iPXE script if not supplying script in entirety
	CustomIPXEScriptURL string `json:"customIPXEScriptURL,omitempty"`
	// Either deploy these Device IDs or fail
	ForceDeviceIds []int32 `json:"forceDeviceIds,omitempty"`
	ProductId      int32   `json:"productId"`
	Options        []int32 `json:"options,omitempty"`
	// must be one of ['monthly', 'quarterly', 'semi-annually', 'annually', 'biennial', 'triennial', 'hourly']
	BillingPeriod string `json:"billingPeriod,omitempty"`
	Quantity      int32  `json:"quantity,omitempty"`
	// ID of SSH Key to use for deployment
	PublicSshKeyId int32 `json:"publicSshKeyId,omitempty"`
	// Operating System's Name or ID
	OperatingSystem string   `json:"operatingSystem"`
	AdditionalNotes []string `json:"additionalNotes,omitempty"`
}