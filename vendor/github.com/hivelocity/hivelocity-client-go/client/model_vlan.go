/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Vlan struct {
	// Unique ID of the VLAN.
	VlanId int32 `json:"vlanId,omitempty"`
	// Unique IDs of ports or bonds.
	PortIds []int32 `json:"portIds,omitempty"`
	// If true, VLAN is configured in Q-in-Q Mode. Automation is not currently supported on Q-in-Q VLANs.
	QInQ bool `json:"qInQ,omitempty"`
	// For example: `NYC1`.
	FacilityCode string `json:"facilityCode,omitempty"`
	// If `public`, this VLAN can have IPs assigned to become reachable from the internet. If `private`, this VLAN can not have IPs assigned and will never be reachabled from the internet. All VLANs are subject to traffic billing and overages, with the exception of private VLAN traffic on unbonded Devices.
	Type_ string `json:"type,omitempty"`
	// Unique IDs of IP Assignments.
	IpIds []int32 `json:"ipIds,omitempty"`
	// If true, VLAN can be automated via API. If false, contact support to enable automation.
	Automated bool `json:"automated,omitempty"`
	// The VLAN Tag id from the switch. Use this value when configuring your OS interfaces to use the VLAN.
	VlanTag int32 `json:"vlanTag,omitempty"`
}