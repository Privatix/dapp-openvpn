<#
.SYNOPSIS
    Creates scheduled task to re-enable Internet Connection Sharing and restart OpenVPN server.
.DESCRIPTION
    There is a known bug with ICS (Internet Connection Sharing), where it must be re-enabled after restart. 
    This script will create scheduled task for this operation.

.PARAMETER scriptPath
    Path to reenable-nat.ps1 script.

.PARAMETER TAPdeviceAddress
    Unique identifier of TAP device. It is identified by "PnPDeviceID" (Get-NetAdapter) and same as "Device instance path" in device manager.

.EXAMPLE
    In cmd.exe with admin priveledges run following command: 
    powershell.exe -executionpolicy bypass -noprofile -file `
    "C:\Program Files\Privatix\agent\product\73e17130-2a1d-4f7d-97a8-93a9aaa6f10d\bin\new-startupTask.ps1" `
    -scriptPath "C:\Program Files\Privatix\agent\product\73e17130-2a1d-4f7d-97a8-93a9aaa6f10d\bin\reenable-nat.ps1" `
    -TAPdeviceAddress "ROOT\NET\0000"

    Description
    -----------
    Creates scheduled task for ICS re-enabling.

#>
Param(
    [ValidateScript( {Test-Path $_})]
    [string]$scriptPath,
    [string]$TAPdeviceAddress
)
#Requires -RunAsAdministrator

$TaskTrigger = New-ScheduledTaskTrigger -AtStartup
$TaskTrigger.Delay = 'PT2M'
$TaskTrigger.ExecutionTimeLimit = 'PT2M'
$PowershellPath = (Get-Command "powershell.exe").Source
$Action = New-ScheduledTaskAction -Execute "$PowershellPath" -Argument "-executionpolicy bypass -noprofile -file `"$scriptPath`" -TAPdeviceAddress `"$TAPdeviceAddress`""
$Settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -RunOnlyIfNetworkAvailable
$User = "NT AUTHORITY\System"
Register-ScheduledTask -TaskName "Privatix re-enable ICS" -Trigger $TaskTrigger -User $User -Action $Action -Settings $Settings -RunLevel Highest -Force -Description "Internet connection sharing must be re-enabled after computer restart, due to Microsoft bug." | Out-Null