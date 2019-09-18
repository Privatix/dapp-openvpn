<#
.SYNOPSIS
    Set IP forwarding for OpenVPN TAP device.
.DESCRIPTION
    Set IP forwarding for OpenVPN TAP device. This operation allows OpenVPN client to browse internet.
    NOTE: internet facing adapter is considered as adapter with lowest route metric. This should be valid in most configurations, but not all.

.PARAMETER TAPdeviceAddress
    Unique identifier of TAP device. It is identified by "PnPDeviceID" (Get-NetAdapter) and same as "Device instance path" in device manager.

.PARAMETER Enabled
    If specified - enable ICS. If ommited - disable.

.PARAMETER Force
    If specified - force to set ICS, even if ICS is already set.

.EXAMPLE
    .\set-nat.ps1 -TAPdeviceAddress 'ROOT\NET\0002' -Enabled

    Description
    -----------
    Enables IP routing. Configures IP forwarding on internet adapter and TAP adapter. Starts "SharedAccess" and "RemoteAccess" services.

.EXAMPLE
    .\set-nat.ps1 -TAPdeviceAddress 'ROOT\NET\0002'

    Description
    -----------
    Disables internet sharing for VPN adapter with specified TAPdeviceAddress
#>

#Requires -Version 3.0 -Modules NetAdapter
#Requires -RunAsAdministrator

param (
    [Parameter(Mandatory)]
    [ValidateNotNullOrEmpty()]
    [string]$TAPdeviceAddress,
    [switch]$Enabled,
    [switch]$Force
)

<#
.SYNOPSIS
    Set internet connection sharing between two adapters
.DESCRIPTION
    Set internet connection sharing between two adapters

.PARAMETER InetAdapterName
    Name of adapter with internet connection

.PARAMETER VPNAdapterName
    Name of TAP adapter, that should get connection to internet

.EXAMPLE
    Set-InternetConnectionSharing -InetAdapterName 'Ethernet' -VPNAdapterName 'Privatix VPN Server'

    Description
    -----------
    Enables internet sharing on "Ethernet" adapter, giving "Privatix VPN Server" adapter to use its internet connection.

#>
function Set-InternetConnectionSharing {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)]
        [ValidateScript(
            { if ((Get-NetAdapter -Name $_ -ErrorAction SilentlyContinue -OutVariable inetAdapter) -and (($inetAdapter).Status -notin @('Disabled', 'Not Present') ))
                { $true }
                else {
                    throw "`"$_`" adapter not exists or disabled"
                }
            }
        )]
        $InetAdapterName,
        [Parameter(Mandatory)]
        [ValidateScript(
            { if ((Get-NetAdapter -Name $_ -ErrorAction SilentlyContinue -OutVariable inetAdapter) -and (($inetAdapter).Status -notin @('Disabled', 'Not Present') ))
                { $true }
                else {
                    throw "`"$_`" adapter not exists or disabled"
                }
            }
        )]
        $VPNAdapterName
    )

    begin {
        $ns = $null

        try {
            # Create a NetSharingManager object
            $ns = New-Object -ComObject HNetCfg.HNetShare
        }
        catch {
            # Register the HNetCfg library (once)
            regsvr32 /s hnetcfg.dll

            # Create a NetSharingManager object
            $ns = New-Object -ComObject HNetCfg.HNetShare
        }

    }

    process {
        # Get internet connected adapter internet connection sharing configuration
        try{
            $InetConn = $ns.EnumEveryConnection | Where-Object { $ns.NetConnectionProps.Invoke($_).Name -eq $InetAdapterName }
            $InetSharingConf = $ns.INetSharingConfigurationForINetConnection.Invoke($InetConn)
        } catch {
            throw "Failed to get internet connected adapter ICS configuration  adapter:`"$InetAdapterName`". Original exception: $($Error[0].exception)"
        }
        # Get VPN server adapter internet connection sharing configuration
        try{
            $VPNConn = $ns.EnumEveryConnection | Where-Object { $ns.NetConnectionProps.Invoke($_).Name -eq $VPNAdapterName }
            $VPNSharingConf = $ns.INetSharingConfigurationForINetConnection.Invoke($VPNConn)
        } catch {
            throw "Failed to get VPN server adapter adapter ICS configuration  adapter:`"$VPNAdapterName`". Original exception: $($Error[0].exception)"
        }

        try {
            $InetSharingConf.EnableSharing(0)
            $VPNSharingConf.EnableSharing(1)
        }
        catch { throw "Failed to enable internet sharing for public adapter `"$InetAdapterName`" and VPN adapter `"$VPNAdapterName`". Original exception: $($Error[0].exception)" }

    }

    end {
        [System.Runtime.Interopservices.Marshal]::ReleaseComObject($ns) | Out-Null
    }
}

<#
.SYNOPSIS
    Reset all previously set ICS
#>
function Reset-InternetConnectionSharing {
    try {
        Get-WmiObject -Namespace "ROOT\Microsoft\HomeNet" -Class "HNet_ConnectionProperties" `
        | Set-WmiInstance -Arguments @{IsIcsPrivate = "false"; IsIcsPublic = "false"}
    } catch { 
        Write-Warning -Message "Failed to reset ICS. Original exception $($error[0].exception)"
    }
    if ((Get-Service -Name "SharedAccess").Status -eq "Running") {
        Restart-Service -Name "SharedAccess"
    }
}
<#
.SYNOPSIS
    Checks if ICS already configured
#>
function Test-ICSconfigured {
    $ICSSet = Get-WmiObject -Namespace "ROOT\Microsoft\HomeNet" -Class "HNet_ConnectionProperties" `
        | Where-Object {$_.IsIcsPrivate -eq $true -or $_.IsIcsPublic -eq $true}
    if ($ICSSet) {return $true} else {return $false}    
}

# Find Internet connected adapter assuming it has lowest metric
$minRouteMetric = (Get-NetRoute | Measure-Object RouteMetric -Minimum).Minimum
$ifIndex = (Get-NetRoute | Where-Object { $_.RouteMetric -eq $minRouteMetric }).ifIndex | Select-Object -Unique
$InetAdapterName = (Get-NetAdapter -Physical | Where-Object { $_.ifIndex -in $ifIndex }).Name
# Find VPN server adapter by TAPdeviceAddress
$VPNAdapterName = (Get-NetAdapter | Where-Object { $_.PnPDeviceID -eq $TAPdeviceAddress }).Name

if ($PSBoundParameters.ContainsKey('Enabled')) {
    # enable routing in registry
    $registryPath = "HKLM:\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters"
    Get-ItemProperty -Path $registryPath -Name "IPEnableRouter" | Set-ItemProperty -Name "IPEnableRouter" -Value 1

    #start windows services
    Get-Service -Name "SharedAccess" | Set-Service -StartupType Automatic | Start-Service
    Get-Service -Name "RemoteAccess" | Set-Service -StartupType Automatic | Start-Service
    
    if (Test-ICSconfigured) {
        if ($PSBoundParameters.ContainsKey('Force')) {
            Reset-InternetConnectionSharing
        } else {
            throw "Connection sharing already enabled on adapter `"$InetAdapterName`". Please, disable it first or use Force flag"
        }
    }
  
    Set-InternetConnectionSharing -InetAdapterName $InetAdapterName -VPNAdapterName $VPNAdapterName
}
else {
    Reset-InternetConnectionSharing
    Get-ItemProperty -Path $registryPath -Name "IPEnableRouter" | Set-ItemProperty -Name "IPEnableRouter" -Value 0
}
