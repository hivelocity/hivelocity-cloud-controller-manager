/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type BareMetalDeviceUpdate struct {
	// A FQDN for the device. For example: `example.hivelocity.net`
	Hostname string `json:"hostname"`
	// The unique ID of an Ignition File for FlatcarOS provisions.
	IgnitionId int32    `json:"ignitionId,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	// URL to download custom iPXE script if not specifying contents in entirety. If both script URL and contents are not specified, the last iPXE script contents are used if OS selection requires an  iPXE script.
	CustomIPXEScriptURL string `json:"customIPXEScriptURL,omitempty"`
	PublicSshKeyId      int32  `json:"publicSshKeyId,omitempty"`
	// A Cloud-Init script or a post-install script (Bash for Linux or Powershell for Windows).
	Script string `json:"script,omitempty"`
	// Contents of iPXE script if not specifying URL. If both script URL and contents are not specified, the last iPXE script contents are used if OS selection requires an iPXE script.
	CustomIPXEScriptContents string `json:"customIPXEScriptContents,omitempty"`
	// The name of the Operating System to provision on this device. Must match name of an operating system product option.
	OsName string `json:"osName"`
}