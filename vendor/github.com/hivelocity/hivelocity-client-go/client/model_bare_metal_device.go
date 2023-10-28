/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type BareMetalDevice struct {
	// The user specified hostname for the device. Note: If the hostname is changed in the portal or on the device itself this value may not reflect the actual hostname on the device.
	Hostname string `json:"hostname,omitempty"`
	// The first assigned public IP for accessing this device.
	PrimaryIp string `json:"primaryIp,omitempty"`
	// User specified values.
	Tags []string `json:"tags,omitempty"`
	// URL of custom iPXE script used to provision device
	CustomIPXEScriptURL string `json:"customIPXEScriptURL,omitempty"`
	// A facility code. For example `NYC1`.
	LocationName string `json:"locationName,omitempty"`
	// The unique ID of the service associated with this device.
	ServiceId int32 `json:"serviceId,omitempty"`
	// The unique ID of the device.
	DeviceId int32 `json:"deviceId,omitempty"`
	// The name of the product associated with this device.
	ProductName string `json:"productName,omitempty"`
	VlanId      int32  `json:"vlanId,omitempty"`
	// This device's service's billing period.
	Period         string `json:"period,omitempty"`
	PublicSshKeyId int32  `json:"publicSshKeyId,omitempty"`
	// The post-install/cloud-init script used during this device's last provisioning.
	Script string `json:"script,omitempty"`
	// ON|OFF
	PowerStatus string `json:"powerStatus,omitempty"`
	// Contents of custom iPXE used to provision device
	CustomIPXEScriptContents string `json:"customIPXEScriptContents,omitempty"`
	// The unique ID of the order for this device.
	OrderId int32 `json:"orderId,omitempty"`
	// The name of the operating system currently installed on this device. Note: If you manually reload your own OS over IPMI this value may not reflect the OS currently on your device.
	OsName string `json:"osName,omitempty"`
	// The unique ID of the product associated with this device.
	ProductId int32 `json:"productId,omitempty"`
}
