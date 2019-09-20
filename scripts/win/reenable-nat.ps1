<#
.SYNOPSIS
    Re-enable IP forwarding for OpenVPN TAP device.
    Restart OpenVPN services.
.DESCRIPTION
    There is a known bug with ICS (Internet Connetcion Sharing), where it must be re-enabled after restart. 
    This script will re-enable ICS and restart Privatix OpenVPN services.

.PARAMETER TAPdeviceAddress
    Unique identifier of TAP device. It is identified by "PnPDeviceID" (Get-NetAdapter) and same as "Device instance path" in device manager.

.EXAMPLE
    .\reenable-nat.ps1 -TAPdeviceAddress 'ROOT\NET\0002'

    Description
    -----------
    Re-enables ICS. Restarts Privatix OpenVPN server.

#>
param(
    [string]$TAPdeviceAddress
)
$ScriptPath = Join-Path $PSScriptRoot -ChildPath "set-nat.ps1"
# Disable ICS
try {
    .$ScriptPath -TAPdeviceAddress $TAPdeviceAddress    
}
catch {
    Write-Warning "Failed to disable Internet connetcion sharing for device: $TAPdeviceAddress"        
}
# Enable ICS
.$ScriptPath -TAPdeviceAddress $TAPdeviceAddress -Enabled -Force

Get-Service "Privatix_OpenVPN*" | Restart-Service -Force

