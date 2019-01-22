<#
.SYNOPSIS
    Makes replacements in Gopkg.toml.template
.DESCRIPTION
    rules: 
    develop         ->  develop
    master          ->  master
    feature/name    ->  develop
    release/name    ->  release/name
    hotfix/name     ->  master


.PARAMETER templateFileName
    Filename of template file

.PARAMETER PROJECT_PATH
    Absolute path to project.

.EXAMPLE
    tompl.ps1 -templateFileName "Gopkg.toml.template" -PROJECT_PATH "C:\Users\tester\go\src\github.com\privatix\dapp-installer"

    Description
    -----------
    replaces branch according to rules (above).

#>
param(
    [ValidateNotNullOrEmpty()]
    [string]$templateFileName,
    [ValidateScript( {Test-Path $_ -PathType "Container"})]
    [string]$PROJECT_PATH
)

$ErrorActionPreference = "Stop"

# get current branch name
$currentBranch = Invoke-Expression "git.exe --git-dir=$PROJECT_PATH\.git --work-tree=$PROJECT_PATH rev-parse --abbrev-ref HEAD"
if ($currentBranch -notmatch "^(?!@$|build-|.*([.]\.|@\{|\\))[^\000-\037\177 ~^:?*[]+[^\000-\037\177 ~^:?*[]+(?<!\.lock|[.])$") {exit 1}

# set replacement by default
$replacement = $currentBranch

# if starts with "feature", then replacement=develop
if ($currentBranch -like "feature*") {
    $replacement = "develop"
}

# if starts with "hotfix", then replacement=master
if ($currentBranch -like "hotfix*") {
    $replacement = "master"
}

# replace "%BRANCH_NAME%" by replacement in the given file
$TomlFilePath = Join-Path -Path $PROJECT_PATH -ChildPath $templateFileName -Resolve
$NewContent = (Get-Content -Path $TomlFilePath | Out-String) -replace "%BRANCH_NAME%", $replacement
$NewContent | Out-File $TomlFilePath -Force -Encoding ascii
