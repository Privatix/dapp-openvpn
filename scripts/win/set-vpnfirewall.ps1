<#
.SYNOPSIS
    Add/Remove Windows firewall rule for OpenVPN server.
.DESCRIPTION
    Add/Remove Windows firewall rule for OpenVPN server.

.PARAMETER Create
    Create firewall rule

.PARAMETER Remove
    Remove firewall rule

.PARAMETER ServiceName
    Name of service. Used for consitent naming between windows service and firewall rule.

.PARAMETER ProgramPath
    Path to binary

.PARAMETER Port
    Port number to allow through firewall

.PARAMETER Protocol
    Protocol. can be udp or tcp

.EXAMPLE
    .\set-vpnfirewall.ps1 -Create -ServiceName "Privatix_openvpn_7b1f782b82f83be7f7eb024def947bc214fa79a3" -ProgramPath "C:\Program Files\Privatix\Agent\73e17130-2a1d-4f7d-97a8-93a9aaa6f10d\bin\openvpn\openvpn.exe" -Port 443 -Protocol tcp

    Description
    -----------
    Allow VPN connection through windows firewall

.EXAMPLE
    .\set-vpnfirewall.ps1 -Remove -ServiceName "Privatix_openvpn_7b1f782b82f83be7f7eb024def947bc214fa79a3" 

    Description
    -----------
    Removes firewall rule for OpenVPN
#>
[cmdletbinding(
    DefaultParameterSetName = 'Create'
)]
param (
    [Parameter(ParameterSetName = "Create", Mandatory = $true)]
    [switch]$Create,
    [Parameter(ParameterSetName = "Remove", Mandatory = $true)]
    [switch]$Remove,
    [Parameter(ParameterSetName = "Create")]
    [Parameter(ParameterSetName = "Remove")]
    [ValidateScript( { Get-Service -Name $_ })]
    [string]$ServiceName,
    [Parameter(ParameterSetName = "Create")]
    [ValidateScript( { Test-Path $_ })]
    [string]$ProgramPath,
    [Parameter(ParameterSetName = "Create")]
    [ValidateRange(0, 65535)] 
    [int]$Port,
    [Parameter(ParameterSetName = "Create")]
    [ValidateSet('tcp', 'udp')]
    [string]$Protocol = 'tcp'
)
if ($PSBoundParameters.ContainsKey('Create')) {
    # Allow inbound connection to OpenVPN server
    New-NetFirewallRule -PolicyStore PersistentStore -Name $ServiceName -DisplayName "Privatix OpenVPN server" `
        -Description "Inbound rule for Privatix OpenVPN server" -Group "Privatix OpenVPN server" -Enabled True -Profile Any `
        -Action Allow -Direction Inbound -LocalPort $Port -Protocol $Protocol -Program $ProgramPath | Out-Null
    # Block connection to LAN
    $LANsubnets = @("10.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12")
    foreach ($LANsubnet in $LANsubnets) {
        try {
            $i++
            New-NetFirewallRule -PolicyStore PersistentStore -Name "Privatix OpenVPN server block LAN access $i" `
                -DisplayName "Privatix OpenVPN server block $LANsubnet" `
                -Description "Outbound rule for Privatix OpenVPN server that block access to LAN" `
                -Group "Privatix OpenVPN server" -Enabled True -Profile Any `
                -Action Block -Direction Outbound -Program $ProgramPath -Protocol Any -RemoteAddress $LANsubnet | Out-Null
        } 
        catch {
            Write-Error "Failed to create firewall rule. Original exception: $($error[0].exception)"
            exit 1
        }
    }

} 
if ($PSBoundParameters.ContainsKey('Remove')) {
    Remove-NetFirewallRule -Group "Privatix OpenVPN server" 
}