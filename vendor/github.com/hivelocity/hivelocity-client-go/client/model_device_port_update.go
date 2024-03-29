/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DevicePortUpdate struct {
	// The unique ID of the port.
	PortId  int32 `json:"portId"`
	Enabled bool  `json:"enabled,omitempty"`
	// IP Assignments IDs currently routed to this port.
	IpAssignments []int32 `json:"ipAssignments"`
	// The unique ID of the port's native vlan, if applicable. Send value `0` to remove the native vlan from this port.
	NativeVlanId int32 `json:"nativeVlanId,omitempty"`
}
