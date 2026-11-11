#!/usr/bin/env pwsh
#
# linkwarden-cli installer for Windows. Works in Windows PowerShell 5.1 and
# PowerShell 7+.
#
# Finds the release for this CPU on GitHub, checks it against the release's
# checksums.txt and installs linkwarden.exe into
# %LOCALAPPDATA%\Programs\linkwarden (no admin needed), or Program Files if you
# pick it from an administrator terminal, plus PowerShell tab completions. If
# the folder isn't on your PATH, it offers to add it.
#
# Usage:
#   irm https://raw.githubusercontent.com/CerealKiller97/linkwarden-cli/main/install.ps1 | iex
#   & ([scriptblock]::Create((irm https://raw.githubusercontent.com/CerealKiller97/linkwarden-cli/main/install.ps1))) -Yes
#
# Run with -Help for the options.
#
# This file is plain ASCII: Windows PowerShell 5.1 reads a .ps1 without a BOM
# in the ANSI code page, so every other character is spelled as a \u escape or
# a code point.
#
# Everything happens inside Install-LinkwardenCli, called at the bottom, so a
# download cut off halfway runs nothing, and `irm | iex` leaves no helper
# functions or preference changes behind in your session.

# Write-Host is the point: it prints to the terminal without feeding the
# pipeline, which `irm | iex` would echo. The helpers are private to this
# script, so they don't take -WhatIf.
[Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSAvoidUsingWriteHost', '')]
[Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSUseShouldProcessForStateChangingFunctions', '')]
param(
  [string]$Version = $env:LINKWARDEN_VERSION,
  [string]$Dir = $env:LINKWARDEN_INSTALL_DIR,
  [switch]$Yes,
  [switch]$Force,
  [switch]$NoCompletions,
  [switch]$NoModifyPath,
  [switch]$Plain,
  [switch]$Quiet,
  [switch]$Uninstall,
  [switch]$Help
)

function Install-LinkwardenCli {
  param([hashtable]$Opt)

  $ErrorActionPreference = 'Stop'
  $ProgressPreference = 'SilentlyContinue'
  Set-StrictMode -Version 2.0

  $Repo = 'CerealKiller97/linkwarden-cli'
  $Bin = 'linkwarden'
  $Exe = 'linkwarden.exe'
  $Steps = 6
  $ProfileMarker = '# Added by the linkwarden-cli installer'
  $UserDir = "$([Environment]::GetFolderPath('LocalApplicationData'))\Programs\linkwarden"
  $SystemDir = "$([Environment]::GetFolderPath('ProgramFiles'))\linkwarden"
  $RawUrl = "https://raw.githubusercontent.com/$Repo/main/install.ps1"

  # Install state. A hashtable, so the helpers below can change it: assigning
  # a plain variable inside a nested function only creates a local copy.
  $S = @{
    ReleasesUrl  = $(if ($env:LINKWARDEN_RELEASES_URL) { $env:LINKWARDEN_RELEASES_URL.TrimEnd('/') } else { "https://github.com/$Repo/releases" })
    Fancy        = $false
    Emoji        = $false
    Interactive  = $false
    Admin        = $false
    TmpDir       = $null
    Client       = $null
    Arch         = $null
    Version      = $null
    Asset        = $null
    InstallDir   = $null
    CursorHidden = $false
    Done         = $false
    Failed       = $false
  }

  $esc = [char]27
  $C = @{ Reset = ''; Bold = ''; Dim = ''; Italic = ''; Red = ''; Green = ''; Yellow = ''; Cyan = ''; White = '' }

  function Get-Glyph([int[]]$CodePoints) {
    -join ($CodePoints | ForEach-Object { [char]::ConvertFromUtf32($_) })
  }

  # Icons print only where the terminal has an emoji font (Windows Terminal,
  # VS Code). The classic console would show empty boxes.
  function Icon([string]$Name) {
    if (-not $S.Emoji) {
      return ''
    }
    $points = @{
      Bookmark = 0x1F516; Package = 0x1F4E6; Key = 0x1F511; Party = 0x1F389; Rocket = 0x1F680
      Wave = 0x1F44B; Lock = 0x1F512; House = 0x1F3E0; Croissant = 0x1F950; Window = 0x1FA9F
      Brain = 0x1F9E0; Shell = 0x1F41A; Muscle = 0x1F4AA; Up = 0x2B06
    }
    $glyph = Get-Glyph $points[$Name]
    if ($Name -eq 'Up') {
      $glyph += [char]0xFE0F
    }
    return "$glyph "
  }

  $G = @{
    Ok = 'OK'; Fail = 'x'; Warn = '!'; Step = '>'; Arrow = '->'; Ask = '?'
    Spin = @('|', '/', '-', '\'); BarDone = '='; BarTodo = '-'; Marker = '>'
    Confetti = @('*', '+', '.', 'o', "'")
  }

  function Write-Out([string]$Text) {
    Write-Host -NoNewline $Text
  }

  # Warnings and errors. On a console, Write-Host: Windows PowerShell 5.1
  # encodes [Console]::Error in the OEM code page, which turns every icon into
  # a "?". Only when stderr is redirected does the text go to the stream.
  function Write-Err([string]$Text) {
    $toStream = $true
    try {
      $toStream = [Console]::IsErrorRedirected
    } catch {
      $toStream = $true
    }
    if ($toStream) {
      [Console]::Error.WriteLine($Text)
    } else {
      Write-Host $Text
    }
  }

  function Stop-Install([string]$Message) {
    throw [System.Exception]::new($Message)
  }

  function Write-Warn([string]$Message) {
    Write-Err ("  {0}{1}{2}  {3}" -f "$($C.Yellow)$($C.Bold)", $G.Warn, $C.Reset, $Message)
  }

  function Write-Info([string]$Message) {
    if (-not $Opt.Quiet) {
      Write-Host "  $Message"
    }
  }

  function Write-Ok([string]$Message) {
    if (-not $Opt.Quiet) {
      Write-Host ("  {0}{1}{2} {3}" -f "$($C.Green)$($C.Bold)", $G.Ok, $C.Reset, $Message)
    }
  }

  function Write-Quip([string]$Message) {
    if (-not $Opt.Quiet) {
      Write-Host ("  {0}{1}{2}" -f "$($C.Dim)$($C.Italic)", $Message, $C.Reset)
    }
  }

  function Wait-Fancy([int]$Milliseconds) {
    if ($S.Fancy) {
      Start-Sleep -Milliseconds $Milliseconds
    }
  }

  function Hide-Cursor {
    if ($S.Fancy) {
      Write-Out "$esc[?25l"
      $S.CursorHidden = $true
    }
  }

  function Show-Cursor {
    if ($S.CursorHidden) {
      Write-Out "$esc[?25h"
      $S.CursorHidden = $false
    }
  }

  function Clear-Line {
    if ($S.Fancy) {
      Write-Out "`r$esc[K"
    }
  }

  # Replaces a leading home folder with ~ for display.
  function Format-Path([string]$Path) {
    $homeDir = $HOME.TrimEnd('\')
    if ($Path -ieq $homeDir) {
      return '~'
    }
    if ($Path.StartsWith("$homeDir\", [StringComparison]::OrdinalIgnoreCase)) {
      return '~' + $Path.Substring($homeDir.Length)
    }
    return $Path
  }

  # Expands ~ and %VARS% and makes the path absolute, relative to the current
  # PowerShell location.
  function Resolve-InstallPath([string]$Path) {
    $p = [Environment]::ExpandEnvironmentVariables($Path.Trim().Trim('"'))
    if ($p -eq '~') {
      $p = $HOME
    } elseif ($p.StartsWith('~\') -or $p.StartsWith('~/')) {
      $p = Join-Path $HOME $p.Substring(2)
    }
    $p = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath($p)
    return $p.TrimEnd('\')
  }

  function Show-Usage {
    Write-Host @"
linkwarden-cli installer for Windows

Usage:
  irm $RawUrl | iex
  & ([scriptblock]::Create((irm $RawUrl))) [options]
  .\install.ps1 [options]

Options:
  -Version <x.y.z>  Install this version instead of the latest.
  -Dir <path>       Install into <path>. Default: $(Format-Path $UserDir).
  -Yes              Don't ask anything; take the defaults. Also on when there's no
                    terminal or CI is set. Never edits your PATH or PowerShell profile.
  -Force            Reinstall even when that version is already installed.
  -NoCompletions    Don't install PowerShell tab completions (installed by default).
  -NoModifyPath     Don't offer to add the install folder to your PATH.
  -Plain            No colours or animations. Also on with NO_COLOR.
  -Quiet            Only print errors and the final result.
  -Uninstall        Remove linkwarden, its completions, and what the installer added
                    to your PATH and PowerShell profile.
  -Help             Show this help.

Environment:
  LINKWARDEN_VERSION       Same as -Version.
  LINKWARDEN_INSTALL_DIR   Same as -Dir.
  LINKWARDEN_RELEASES_URL  Base URL of the releases (default: $($S.ReleasesUrl)).
"@
  }

  # Decides whether to use colours, animations, emoji and prompts.
  function Initialize-Terminal {
    $vt = $false
    try {
      $vt = [bool]$Host.UI.SupportsVirtualTerminal
    } catch {
      $vt = $false
    }
    $redirected = $false
    try {
      $redirected = [Console]::IsOutputRedirected
    } catch {
      $redirected = $true
    }
    if ($vt -and -not $redirected -and -not $env:NO_COLOR -and -not $Opt.Plain -and -not $Opt.Quiet) {
      $S.Fancy = $true
    }
    if ($S.Fancy -and ($env:WT_SESSION -or $env:TERM_PROGRAM)) {
      $S.Emoji = $true
    }
    $inputRedirected = $true
    try {
      $inputRedirected = [Console]::IsInputRedirected
    } catch {
      $inputRedirected = $true
    }
    if (-not $Opt.Yes -and -not $env:CI -and [Environment]::UserInteractive -and -not $inputRedirected) {
      $S.Interactive = $true
    }
    if ($S.Fancy) {
      $C.Reset = "$esc[0m"; $C.Bold = "$esc[1m"; $C.Dim = "$esc[2m"; $C.Italic = "$esc[3m"
      $C.Red = "$esc[31m"; $C.Green = "$esc[32m"; $C.Yellow = "$esc[33m"; $C.Cyan = "$esc[36m"
      $C.White = "$esc[97m"
    }
    if ($S.Emoji) {
      $G.Ok = [regex]::Unescape('\u2714'); $G.Fail = [regex]::Unescape('\u2716')
      $G.Warn = [regex]::Unescape('\u26A0'); $G.Step = [regex]::Unescape('\u25C6')
      $G.Arrow = [regex]::Unescape('\u2192')
      $G.Spin = @([regex]::Unescape('\u280B\u2819\u2839\u2838\u283C\u2834\u2826\u2827\u2807\u280F').ToCharArray() | ForEach-Object { "$_" })
      $G.BarDone = [regex]::Unescape('\u2501'); $G.BarTodo = [regex]::Unescape('\u2500')
      $G.Marker = Get-Glyph 0x1F516
      $G.Confetti = @('*', '+', [regex]::Unescape('\u2726'), [regex]::Unescape('\u2727'),
        [regex]::Unescape('\u2022'), [regex]::Unescape('\u00B7'), [regex]::Unescape('\u00B0'),
        [regex]::Unescape('\u2736'), [regex]::Unescape('\u22C6'))
    }
    try {
      $identity = [Security.Principal.WindowsIdentity]::GetCurrent()
      $S.Admin = ([Security.Principal.WindowsPrincipal]$identity).IsInRole(
        [Security.Principal.WindowsBuiltInRole]::Administrator)
    } catch {
      $S.Admin = $false
    }
  }

  # A point on the sky -> indigo -> fuchsia gradient, as a 24-bit colour code.
  function Get-Gradient([int]$Pos, [int]$Total) {
    if (-not $S.Fancy) {
      return ''
    }
    if ($Total -lt 2) {
      $Total = 2
    }
    $t = [int][Math]::Truncate($Pos * 1000 / ($Total - 1))
    if ($t -le 500) {
      $r = 56 + [int][Math]::Truncate((129 - 56) * $t / 500)
      $g = 189 + [int][Math]::Truncate((140 - 189) * $t / 500)
      $b = 248
    } else {
      $t -= 500
      $r = 129 + [int][Math]::Truncate((232 - 129) * $t / 500)
      $g = 140 + [int][Math]::Truncate((121 - 140) * $t / 500)
      $b = 248 + [int][Math]::Truncate((249 - 248) * $t / 500)
    }
    return "$esc[38;2;$r;$g;${b}m"
  }

  $LogoTop = @('\u257B  ', '\u257B', '\u250F\u2513\u257B', '\u257B\u250F ', '\u257B \u257B', '\u250F\u2501\u2513', '\u250F\u2501\u2513', '\u257A\u2533\u2513', '\u250F\u2501\u2578', '\u250F\u2513\u257B')
  $LogoMid = @('\u2503  ', '\u2503', '\u2503\u2517\u252B', '\u2523\u253B\u2513', '\u2503\u257B\u2503', '\u2523\u2501\u252B', '\u2523\u2533\u251B', ' \u2503\u2503', '\u2523\u2578 ', '\u2503\u2517\u252B')
  $LogoBot = @('\u2517\u2501\u2578', '\u2579', '\u2579 \u2579', '\u2579 \u2579', '\u2517\u253B\u251B', '\u2579 \u2579', '\u2579\u2517\u2578', '\u257A\u253B\u251B', '\u2517\u2501\u2578', '\u2579 \u2579')

  # One row of the logo, letter by letter along the gradient. The letter at
  # $Highlight is drawn bright white, for the shimmer.
  function Write-LogoRow([string[]]$Row, [int]$Highlight, [bool]$Cli) {
    $out = '  '
    for ($i = 0; $i -lt $Row.Count; $i++) {
      $glyph = [regex]::Unescape($Row[$i])
      if ($i -eq $Highlight) {
        $out += "$($C.Bold)$($C.White)$glyph$($C.Reset)"
      } else {
        $out += "$($C.Bold)$(Get-Gradient $i $Row.Count)$glyph$($C.Reset)"
      }
    }
    if ($Cli) {
      $out += "  $($C.Dim)cli$($C.Reset)"
    }
    Write-Host "$out$esc[K"
  }

  # Prints text one character at a time, like it's being typed.
  function Write-Typed([string]$Text) {
    if (-not $S.Fancy) {
      Write-Host $Text
      return
    }
    foreach ($ch in $Text.ToCharArray()) {
      Write-Out "$ch"
      Start-Sleep -Milliseconds 12
    }
    Write-Host ''
  }

  # The opening: the logo drops in row by row, a shimmer sweeps across it, and
  # the tagline types itself out. The logo needs Windows Terminal's fonts.
  function Show-Intro {
    if ($Opt.Quiet) {
      return
    }
    Write-Host ''
    if (-not $S.Emoji) {
      if ($S.Fancy) {
        Write-Host "  $($C.Bold)$(Get-Gradient 0 2)linkwarden$($C.Reset) $($C.Dim)cli installer$($C.Reset)"
        Write-Host ''
        Write-Out '  '
        Write-Typed 'Your bookmarks, from the comfort of your terminal.'
        Write-Host ''
      } else {
        Write-Host '  linkwarden cli installer'
        Write-Host ''
      }
      return
    }
    Hide-Cursor
    Write-LogoRow -Row $LogoTop -Highlight -1 -Cli $false
    Start-Sleep -Milliseconds 80
    Write-LogoRow -Row $LogoMid -Highlight -1 -Cli $true
    Start-Sleep -Milliseconds 80
    Write-LogoRow -Row $LogoBot -Highlight -1 -Cli $false
    for ($hl = 0; $hl -le 10; $hl++) {
      Write-Out "$esc[3A"
      Write-LogoRow -Row $LogoTop -Highlight $hl -Cli $false
      Write-LogoRow -Row $LogoMid -Highlight $hl -Cli $true
      Write-LogoRow -Row $LogoBot -Highlight $hl -Cli $false
      Start-Sleep -Milliseconds 50
    }
    Write-Host ''
    Write-Out "  $(Icon 'Bookmark')"
    Write-Typed 'Your bookmarks, from the comfort of your terminal.'
    Write-Host ''
  }

  # A step heading such as "> 2/6  Finding the latest release".
  function Write-Step([int]$Number, [string]$Title) {
    if ($Opt.Quiet) {
      return
    }
    if ($S.Fancy) {
      Write-Host ''
      Write-Host ("{0}{1}{2} {3}/{4}{5}  {6}{7}{8}" -f $C.Bold, (Get-Gradient ($Number - 1) $Steps), $G.Step,
        $Number, $Steps, $C.Reset, $C.Bold, $Title, $C.Reset)
    } else {
      Write-Host ''
      Write-Host "[$Number/$Steps] $Title"
    }
  }

  # Waits for a .NET task while a spinner turns, then returns its result.
  # Polling with Start-Sleep keeps Ctrl-C working while it waits.
  function Wait-Task($Task, [string]$Label) {
    $i = 0
    if ($S.Fancy -and $Label) {
      Hide-Cursor
      while (-not $Task.IsCompleted -or $i -lt 8) {
        Write-Out ("`r  {0}{1}{2} {3}" -f $C.Cyan, $G.Spin[$i % $G.Spin.Count], $C.Reset, $Label)
        $i++
        Start-Sleep -Milliseconds 70
      }
      Clear-Line
    } else {
      while (-not $Task.IsCompleted) {
        Start-Sleep -Milliseconds 50
      }
    }
    return $Task.GetAwaiter().GetResult()
  }

  # The innermost message of an error, without PowerShell's
  # "Exception calling ... with 0 argument(s)" wrapping.
  function Get-ErrorText($ErrorRecord) {
    $ex = $ErrorRecord.Exception
    while ($ex.InnerException) {
      $ex = $ex.InnerException
    }
    return $ex.Message
  }

  # Asks a yes/no question. Without a terminal, returns the default.
  function Confirm-Choice([string]$Question, [bool]$Default) {
    if (-not $S.Interactive) {
      return $Default
    }
    $hint = $(if ($Default) { 'Y/n' } else { 'y/N' })
    Show-Cursor
    Write-Out ("  {0}{1}{2} {3} {4}({5}){6} " -f "$($C.Cyan)$($C.Bold)", $G.Ask, $C.Reset, $Question, $C.Dim, $hint, $C.Reset)
    $answer = Read-Host
    if ([string]::IsNullOrWhiteSpace($answer)) {
      return $Default
    }
    return $answer.Trim() -match '^[Yy]'
  }

  # Reads a line of text.
  function Read-Answer([string]$Prompt) {
    Show-Cursor
    Write-Out ("  {0}{1}{2} {3} " -f "$($C.Cyan)$($C.Bold)", $G.Ask, $C.Reset, $Prompt)
    return "$(Read-Host)".Trim()
  }

  # Sets $S.Arch to GoReleaser's name for this machine's CPU.
  function Get-Platform {
    $isCore = $PSVersionTable.PSEdition -eq 'Core'
    if ($isCore -and -not (Get-Variable -Name IsWindows -ValueOnly -ErrorAction SilentlyContinue)) {
      Stop-Install "this installer is for Windows. On macOS and Linux, run:  curl -fsSL https://raw.githubusercontent.com/$Repo/main/install.sh | bash"
    }
    $arch = $null
    try {
      $arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString()
    } catch {
      $arch = $null
    }
    if (-not $arch) {
      $arch = $(if ($env:PROCESSOR_ARCHITEW6432) { $env:PROCESSOR_ARCHITEW6432 } else { $env:PROCESSOR_ARCHITECTURE })
    }
    switch -Regex ($arch) {
      '^(X64|AMD64)$' { $S.Arch = 'amd64' }
      '^(Arm64|ARM64)$' { $S.Arch = 'arm64' }
      default {
        Stop-Install "no prebuilt linkwarden for $arch CPUs. Try: go install github.com/$Repo/cmd/linkwarden@latest"
      }
    }
    # An x64 PowerShell on Windows on ARM may report X64; the native build is faster.
    if ($S.Arch -eq 'amd64' -and "$env:PROCESSOR_IDENTIFIER" -match '\bARM') {
      $S.Arch = 'arm64'
    }
  }

  # Prints what Get-Platform found, with a little commentary.
  function Show-Platform {
    $cpu = $(if ($S.Arch -eq 'arm64') { 'ARM64' } else { 'x86_64' })
    $shell = "PowerShell $($PSVersionTable.PSVersion.Major).$($PSVersionTable.PSVersion.Minor)"
    Wait-Fancy 150
    Write-Info "$(Icon 'Window')$($C.Dim)OS   $($C.Reset) $($C.Bold)Windows$($C.Reset) $($C.Dim)(windows)$($C.Reset)"
    Wait-Fancy 150
    Write-Info "$(Icon 'Brain')$($C.Dim)CPU  $($C.Reset) $($C.Bold)$cpu$($C.Reset) $($C.Dim)($($S.Arch))$($C.Reset)"
    Wait-Fancy 150
    Write-Info "$(Icon 'Shell')$($C.Dim)Shell$($C.Reset) $($C.Bold)$shell$($C.Reset)"
    Wait-Fancy 150
    if ($S.Arch -eq 'arm64') {
      Write-Quip "Windows on ARM. Bold. $(Icon 'Muscle')"
    } else {
      Write-Quip "Windows on x86_64. A trusty workhorse. $(Icon 'Window')"
    }
  }

  function New-HttpClient {
    try {
      # Windows PowerShell 5.1 may still default to TLS 1.0, which GitHub refuses.
      [Net.ServicePointManager]::SecurityProtocol = [Net.ServicePointManager]::SecurityProtocol -bor [Net.SecurityProtocolType]::Tls12
    } catch {
      Write-Verbose 'Could not enable TLS 1.2; carrying on with the defaults.'
    }
    Add-Type -AssemblyName System.Net.Http
    $handler = [System.Net.Http.HttpClientHandler]::new()
    $handler.AllowAutoRedirect = $true
    $client = [System.Net.Http.HttpClient]::new($handler)
    $client.Timeout = [TimeSpan]::FromMinutes(10)
    $client.DefaultRequestHeaders.UserAgent.ParseAdd('linkwarden-cli-installer')
    $S.Client = $client
  }

  # Sets $S.Version and $S.Asset, from -Version or the latest release. The
  # latest version comes from the releases/latest redirect, so no API token or
  # rate limit is involved.
  function Resolve-Version {
    if ($Opt.Version) {
      $S.Version = $Opt.Version
      Write-Ok "Version $($C.Bold)$($S.Version)$($C.Reset), as requested"
    } else {
      $request = [System.Net.Http.HttpRequestMessage]::new([System.Net.Http.HttpMethod]::Head, "$($S.ReleasesUrl)/latest")
      try {
        $response = Wait-Task $S.Client.SendAsync($request) "Asking GitHub what's new..."
      } catch {
        Stop-Install "couldn't find the latest release: $(Get-ErrorText $_)"
      }
      $tag = $response.RequestMessage.RequestUri.AbsolutePath.Split('/')[-1]
      $response.Dispose()
      if ($tag -notmatch '^v\d') {
        Stop-Install "no release found at $($S.ReleasesUrl)/latest"
      }
      $S.Version = $tag.Substring(1)
      Write-Ok "Latest release is $($C.Bold)v$($S.Version)$($C.Reset)"
      Write-Quip "Fresh out of the oven. $(Icon 'Croissant')"
    }
    $S.Asset = "${Bin}_$($S.Version)_windows_$($S.Arch).zip"
  }

  # The version an installed linkwarden.exe reports, without the "v".
  function Get-InstalledVersion([string]$Path) {
    $ErrorActionPreference = 'Continue'
    try {
      $out = (& $Path --version 2>$null | Out-String).Trim()
    } catch {
      return ''
    }
    if (-not $out) {
      return ''
    }
    return ($out -split '\s+')[-1].TrimStart('v')
  }

  # Looks for an installed linkwarden. Returns $false when it's already this
  # version and there's nothing to do.
  function Test-Existing {
    $found = Get-Command $Bin -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
    if (-not $found) {
      return $true
    }
    $existing = Get-InstalledVersion $found.Path
    if ($existing -ne $S.Version) {
      $shownExisting = $(if ($existing) { $existing } else { 'unknown' })
      Write-Info ("{0}Upgrade: {1}{2}{3} {4} {5}{6}{7}{8} {9}({10}){11}" -f (Icon 'Up'), $C.Dim, $shownExisting, $C.Reset,
        $G.Arrow, $C.Bold, $C.Green, $S.Version, $C.Reset, $C.Dim, (Format-Path $found.Path), $C.Reset)
      return $true
    }
    Write-Ok "linkwarden $($S.Version) is already installed at $(Format-Path $found.Path)"
    if ($Opt.Force) {
      return $true
    }
    if (Confirm-Choice 'Reinstall it anyway?' $false) {
      return $true
    }
    Write-Host ''
    Write-Host "  $($C.Green)All good, nothing to do. Happy bookmarking! $(Icon 'Bookmark')$($C.Reset)"
    Write-Host ''
    return $false
  }

  # Exits with an error when the current user can't write to the folder.
  function Assert-Writable([string]$Folder) {
    try {
      New-Item -ItemType Directory -Force -Path $Folder | Out-Null
      $probe = Join-Path $Folder ".linkwarden-install-$PID"
      [IO.File]::WriteAllText($probe, '')
      Remove-Item -LiteralPath $probe -Force
    } catch {
      $hint = $(if ($S.Admin) { 'Pick another folder with -Dir.' } else { 'Run the installer from an administrator terminal, or pick another folder with -Dir.' })
      Stop-Install "can't write to $Folder. $hint"
    }
  }

  # Sets $S.InstallDir: -Dir, or a choice from a small menu, or the default,
  # the per-user folder.
  function Select-InstallDir {
    if ($Opt.Dir) {
      $S.InstallDir = Resolve-InstallPath $Opt.Dir
    } elseif ($S.Interactive) {
      Write-Host "  $($C.Bold)Where should linkwarden live?$($C.Reset)"
      Write-Host ("    {0}1){1} {2,-36} {3}<- just for you, no admin (default){4}" -f $C.Cyan, $C.Reset, (Format-Path $UserDir), $C.Dim, $C.Reset)
      Write-Host ("    {0}2){1} {2,-36} {3}<- everyone on this PC, needs an admin terminal{4}" -f $C.Cyan, $C.Reset, $SystemDir, $C.Dim, $C.Reset)
      Write-Host ("    {0}3){1} somewhere else..." -f $C.Cyan, $C.Reset)
      $choice = Read-Answer 'Pick one [1-3]:'
      if (-not $choice) {
        $choice = '1'
      }
      switch ($choice) {
        '1' { $S.InstallDir = $UserDir }
        '2' {
          if (-not $S.Admin) {
            Stop-Install "$SystemDir needs an administrator terminal. Open one and run the installer again, or pick 1."
          }
          $S.InstallDir = $SystemDir
        }
        '3' {
          $folder = Read-Answer 'Folder:'
          if (-not $folder) {
            Stop-Install 'no folder given'
          }
          $S.InstallDir = Resolve-InstallPath $folder
        }
        default { Stop-Install "'$choice' isn't one of 1, 2 or 3" }
      }
      Hide-Cursor
    } else {
      $S.InstallDir = $UserDir
    }
    Assert-Writable $S.InstallDir
    Write-Ok "Installing to $($C.Bold)$(Format-Path $S.InstallDir)$($C.Reset)"
    Write-Quip "Cozy. $(Icon 'House')"
  }

  # Formats a byte count as e.g. "4.2 MB".
  function Format-Size([double]$Bytes) {
    if ($Bytes -ge 1MB) {
      return '{0:N1} MB' -f ($Bytes / 1MB)
    }
    if ($Bytes -ge 1KB) {
      return '{0:N0} KB' -f ($Bytes / 1KB)
    }
    return "$([int64]$Bytes) B"
  }

  # Draws the download progress bar, with a bookmark running along it.
  function Write-Bar([int64]$DoneBytes, [int64]$Total) {
    $width = 32
    $filled = [int][Math]::Min($width, [Math]::Truncate($DoneBytes * $width / $Total))
    $pct = [int][Math]::Truncate($DoneBytes * 100 / $Total)
    $bar = ''
    for ($i = 0; $i -lt $width; $i++) {
      if ($i -lt $filled) {
        $bar += "$(Get-Gradient $i $width)$($G.BarDone)"
      } elseif ($i -eq $filled) {
        $bar += "$($C.Reset)$($G.Marker)"
        if ($S.Emoji) {
          # The bookmark is two columns wide.
          $i++
        }
      } else {
        $bar += "$($C.Dim)$($G.BarTodo)"
      }
    }
    Write-Out ("`r  {0}{1} {2}{3,3}%{4} {5}{6} / {7}{8}$esc[K" -f $bar, $C.Reset, $C.Bold, $pct, $C.Reset, $C.Dim,
      (Format-Size $DoneBytes), (Format-Size $Total), $C.Reset)
  }

  # Downloads a URL to a file, with a progress bar when there's a terminal.
  function Save-Url([string]$Url, [string]$Destination) {
    $response = Wait-Task $S.Client.GetAsync($Url, [System.Net.Http.HttpCompletionOption]::ResponseHeadersRead) ''
    try {
      if (-not $response.IsSuccessStatusCode) {
        throw [System.Exception]::new("HTTP $([int]$response.StatusCode) $($response.ReasonPhrase)")
      }
      $total = $response.Content.Headers.ContentLength
      $in = Wait-Task $response.Content.ReadAsStreamAsync() ''
      $out = [IO.File]::Create($Destination)
      try {
        $buffer = New-Object byte[] 65536
        [int64]$doneBytes = 0
        $lastDraw = [DateTime]::MinValue
        while (($n = $in.Read($buffer, 0, $buffer.Length)) -gt 0) {
          $out.Write($buffer, 0, $n)
          $doneBytes += $n
          if ($S.Fancy -and $total -and ([DateTime]::UtcNow - $lastDraw).TotalMilliseconds -ge 50) {
            Write-Bar $doneBytes $total
            $lastDraw = [DateTime]::UtcNow
          }
        }
      } finally {
        $out.Dispose()
        $in.Dispose()
      }
      if ($S.Fancy -and $total) {
        Write-Bar $total $total
        Write-Host ''
      }
    } finally {
      $response.Dispose()
    }
  }

  # Downloads the release archive and its checksums.txt into the temp folder.
  function Get-Release {
    $base = "$($S.ReleasesUrl)/download/v$($S.Version)"
    Write-Info "$(Icon 'Package')$($C.Dim)$($S.Asset)$($C.Reset)"
    Hide-Cursor
    try {
      Save-Url "$base/$($S.Asset)" (Join-Path $S.TmpDir $S.Asset)
    } catch {
      Clear-Line
      Stop-Install "couldn't download $base/$($S.Asset) ($(Get-ErrorText $_)). Is v$($S.Version) a real release? See $($S.ReleasesUrl)"
    }
    try {
      $checksums = Wait-Task $S.Client.GetStringAsync("$base/checksums.txt") 'Grabbing checksums.txt...'
    } catch {
      Stop-Install "couldn't download $base/checksums.txt ($(Get-ErrorText $_))"
    }
    [IO.File]::WriteAllText((Join-Path $S.TmpDir 'checksums.txt'), $checksums)
    Write-Ok "Downloaded $(Format-Size (Get-Item -LiteralPath (Join-Path $S.TmpDir $S.Asset)).Length)"
  }

  # Checks the archive against checksums.txt and exits when they differ.
  function Test-Checksum {
    $expected = $null
    foreach ($line in [IO.File]::ReadAllLines((Join-Path $S.TmpDir 'checksums.txt'))) {
      $parts = $line.Trim() -split '\s+'
      if ($parts.Count -ge 2 -and $parts[1].TrimStart('*') -eq $S.Asset) {
        $expected = $parts[0].ToLowerInvariant()
      }
    }
    if (-not $expected) {
      Stop-Install "checksums.txt has no entry for $($S.Asset)"
    }
    $actual = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $S.TmpDir $S.Asset)).Hash.ToLowerInvariant()
    if ($actual -ne $expected) {
      Write-Err ''
      Write-Err "  $($C.Dim)expected $expected$($C.Reset)"
      Write-Err "  $($C.Dim)got      $actual$($C.Reset)"
      Stop-Install "checksum mismatch for $($S.Asset). The download is corrupt or was tampered with; nothing was installed."
    }
    Write-Ok "SHA-256 matches $($C.Dim)$($actual.Substring(0, 12))...$($C.Reset)"
    Write-Quip "Nobody's been messing with your bookmarks. $(Icon 'Lock')"
  }

  # The folder the archive was extracted into, named after the asset.
  function Get-ExtractedDir {
    return Join-Path (Join-Path $S.TmpDir 'extract') ([IO.Path]::GetFileNameWithoutExtension($S.Asset))
  }

  # Extracts the archive and installs linkwarden.exe into $S.InstallDir.
  function Install-Binary {
    $extract = Join-Path $S.TmpDir 'extract'
    try {
      Add-Type -AssemblyName System.IO.Compression.FileSystem
      [IO.Compression.ZipFile]::ExtractToDirectory((Join-Path $S.TmpDir $S.Asset), $extract)
    } catch {
      Stop-Install "couldn't extract $($S.Asset): $(Get-ErrorText $_)"
    }
    $src = Join-Path (Get-ExtractedDir) $Exe
    if (-not (Test-Path -LiteralPath $src)) {
      Stop-Install "$($S.Asset) doesn't contain $Exe"
    }
    $dest = Join-Path $S.InstallDir $Exe
    try {
      Copy-Item -LiteralPath $src -Destination $dest -Force
    } catch {
      Stop-Install "couldn't write $dest. If linkwarden is running, close it and try again. ($(Get-ErrorText $_))"
    }
    Write-Ok "Installed $($C.Bold)$(Format-Path $dest)$($C.Reset)"
  }

  # The profile every PowerShell host of this edition loads.
  function Get-ProfilePath {
    $p = Get-Variable -Name PROFILE -ValueOnly -ErrorAction SilentlyContinue
    if ($p -and $p.PSObject.Properties['CurrentUserAllHosts']) {
      return $p.CurrentUserAllHosts
    }
    $folder = $(if ($PSVersionTable.PSEdition -eq 'Core') { 'PowerShell' } else { 'WindowsPowerShell' })
    return Join-Path (Join-Path ([Environment]::GetFolderPath('MyDocuments')) $folder) 'profile.ps1'
  }

  # Where the completion script goes: next to the binary.
  function Get-CompletionPath([string]$Folder) {
    return Join-Path (Join-Path $Folder 'completions') "$Bin.ps1"
  }

  # Adds a line to the PowerShell profile that loads the completions, after
  # asking. Never edits anything without a yes on the terminal.
  function Add-ProfileLine([string]$CompletionFile) {
    $profilePath = Get-ProfilePath
    $quoted = $CompletionFile.Replace("'", "''")
    $line = "if (Test-Path '$quoted') { . '$quoted' }"
    if ((Test-Path -LiteralPath $profilePath) -and
      (Get-Content -LiteralPath $profilePath -Raw -ErrorAction SilentlyContinue) -like "*$CompletionFile*") {
      Write-Ok "$(Format-Path $profilePath) already loads them"
      return
    }
    $policy = "$(Get-ExecutionPolicy)"
    if ($policy -eq 'Restricted' -or $policy -eq 'AllSigned') {
      Write-Warn "PowerShell's execution policy ($policy) stops your profile from loading them"
      Write-Info "   Allow local scripts with:  $($C.Cyan)Set-ExecutionPolicy -Scope CurrentUser RemoteSigned$($C.Reset)"
      Write-Info "   then add this to $(Format-Path $profilePath):  $($C.Cyan)$line$($C.Reset)"
      return
    }
    if (-not (Confirm-Choice "Load them from your PowerShell profile ($(Format-Path $profilePath))?" $true)) {
      Write-Info "   Load them yourself by adding this to $(Format-Path $profilePath):"
      Write-Info "   $($C.Cyan)$line$($C.Reset)"
      return
    }
    try {
      New-Item -ItemType Directory -Force -Path (Split-Path -Parent $profilePath) | Out-Null
      Add-Content -LiteralPath $profilePath -Value @('', $ProfileMarker, $line) -Encoding UTF8
    } catch {
      Write-Warn "couldn't write $(Format-Path $profilePath): $(Get-ErrorText $_)"
      return
    }
    Write-Ok "Added to $(Format-Path $profilePath). New PowerShell windows get tab completion."
  }

  # Installs PowerShell completions next to the binary. Skipped with
  # -NoCompletions.
  function Install-CompletionScript {
    if ($Opt.NoCompletions) {
      return
    }
    $src = Join-Path (Get-ExtractedDir) "completions\$Bin.ps1"
    if (-not (Test-Path -LiteralPath $src)) {
      # Older releases don't ship a PowerShell script; generate it like
      # scripts/completions.sh does. Any token will do for `completion`.
      $src = Join-Path $S.TmpDir "$Bin.ps1"
      $savedToken = $env:LINKWARDEN_TOKEN
      $ErrorActionPreference = 'Continue'
      try {
        $env:LINKWARDEN_TOKEN = 'completions'
        $completionText = (& (Join-Path $S.InstallDir $Exe) completion powershell 2>$null | Out-String)
        if ($LASTEXITCODE -ne 0 -or -not $completionText.Trim()) {
          return
        }
        [IO.File]::WriteAllText($src, $completionText)
      } catch {
        return
      } finally {
        $env:LINKWARDEN_TOKEN = $savedToken
        $ErrorActionPreference = 'Stop'
      }
    }
    $dest = Get-CompletionPath $S.InstallDir
    try {
      New-Item -ItemType Directory -Force -Path (Split-Path -Parent $dest) | Out-Null
      Copy-Item -LiteralPath $src -Destination $dest -Force
    } catch {
      Write-Warn "Skipped tab completions: couldn't write $dest"
      return
    }
    Write-Ok "Tab completions for PowerShell $($G.Arrow) $(Format-Path $dest)"
    Add-ProfileLine $dest
  }

  # The registry key holding the User or Machine environment.
  function Open-EnvironmentKey([string]$Scope, [bool]$Writable) {
    if ($Scope -eq 'Machine') {
      return [Microsoft.Win32.Registry]::LocalMachine.OpenSubKey('SYSTEM\CurrentControlSet\Control\Session Manager\Environment', $Writable)
    }
    return [Microsoft.Win32.Registry]::CurrentUser.OpenSubKey('Environment', $Writable)
  }

  # The raw PATH of a scope, with %VARS% left unexpanded so rewriting it
  # doesn't freeze them.
  function Get-RawPath([string]$Scope) {
    $key = Open-EnvironmentKey $Scope $false
    if (-not $key) {
      return ''
    }
    try {
      return "$($key.GetValue('Path', '', [Microsoft.Win32.RegistryValueOptions]::DoNotExpandEnvironmentNames))"
    } finally {
      $key.Dispose()
    }
  }

  function Set-RawPath([string]$Scope, [string]$Value) {
    $key = Open-EnvironmentKey $Scope $true
    try {
      $key.SetValue('Path', $Value, [Microsoft.Win32.RegistryValueKind]::ExpandString)
    } finally {
      $key.Dispose()
    }
    # Setting any variable through .NET broadcasts WM_SETTINGCHANGE, so
    # Explorer and new terminals see the new PATH without signing out.
    [Environment]::SetEnvironmentVariable('LINKWARDEN_INSTALLER', '1', $Scope)
    [Environment]::SetEnvironmentVariable('LINKWARDEN_INSTALLER', $null, $Scope)
  }

  function Test-SameFolder([string]$A, [string]$B) {
    return [Environment]::ExpandEnvironmentVariables($A).TrimEnd('\') -ieq $B.TrimEnd('\')
  }

  function Test-OnPath([string]$Folder) {
    foreach ($entry in @($env:Path -split ';') + @((Get-RawPath 'User') -split ';') + @((Get-RawPath 'Machine') -split ';')) {
      if ($entry -and (Test-SameFolder $entry $Folder)) {
        return $true
      }
    }
    return $false
  }

  # When the install folder isn't on PATH, offers to add it: to the user's
  # PATH, or the system's for an admin installing outside their profile.
  function Add-ToPath {
    $folder = $S.InstallDir
    if (Test-OnPath $folder) {
      return
    }
    $insideProfile = $folder.StartsWith("$($HOME.TrimEnd('\'))\", [StringComparison]::OrdinalIgnoreCase)
    $scope = $(if ($S.Admin -and -not $insideProfile) { 'Machine' } else { 'User' })
    $scopeLabel = $(if ($scope -eq 'Machine') { 'the system PATH' } else { 'your PATH' })
    $manual = "[Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', '$scope') + ';$folder', '$scope')"
    Write-Warn "$(Format-Path $folder) isn't on your PATH yet"
    if ($Opt.NoModifyPath -or -not $S.Interactive -or -not (Confirm-Choice "Add it to ${scopeLabel}?" $true)) {
      Write-Info "   Add it yourself with:  $($C.Cyan)$manual$($C.Reset)"
      return
    }
    $raw = (Get-RawPath $scope).TrimEnd(';')
    $new = $(if ($raw) { "$raw;$folder" } else { $folder })
    try {
      Set-RawPath $scope $new
    } catch {
      Write-Warn "couldn't update $scopeLabel`: $(Get-ErrorText $_)"
      Write-Info "   Add it yourself with:  $($C.Cyan)$manual$($C.Reset)"
      return
    }
    $env:Path = "$($env:Path.TrimEnd(';'));$folder"
    if ($Opt.InSession) {
      Write-Ok "Added to $scopeLabel. This window and new terminals can run $Bin now."
    } else {
      Write-Ok "Added to $scopeLabel. Open a new terminal to use $Bin."
    }
  }

  # Confetti falling for a moment, then wiped away.
  function Show-Confetti {
    $rows = 5
    $width = 46
    $lines = New-Object string[] $rows
    for ($r = 0; $r -lt $rows; $r++) {
      $lines[$r] = ''
      Write-Host ''
    }
    $rand = New-Object System.Random
    for ($frame = 0; $frame -lt 18; $frame++) {
      $row = ''
      for ($c = 0; $c -lt $width; $c++) {
        if ($rand.Next(7) -eq 0) {
          $row += "$(Get-Gradient ($rand.Next(10)) 10)$($G.Confetti[$rand.Next($G.Confetti.Count)])"
        } else {
          $row += ' '
        }
      }
      # Shift every row down one and drop the new one in on top.
      for ($r = $rows - 1; $r -gt 0; $r--) {
        $lines[$r] = $lines[$r - 1]
      }
      $lines[0] = "$row$($C.Reset)"
      Write-Out "$esc[${rows}A"
      foreach ($l in $lines) {
        Write-Host "  $l$esc[K"
      }
      Start-Sleep -Milliseconds 60
    }
    Write-Out "$esc[${rows}A"
    for ($r = 0; $r -lt $rows; $r++) {
      Write-Host "$esc[K"
    }
    Write-Out "$esc[${rows}A"
  }

  # The ending: confetti, proof the binary runs, and what to do next.
  function Show-Finale {
    $path = Join-Path $S.InstallDir $Exe
    $installed = Get-InstalledVersion $path
    if (-not $installed) {
      $installed = $S.Version
    }
    if ($Opt.Quiet) {
      Write-Host "linkwarden $installed installed to $path"
      return
    }
    if ($S.Fancy) {
      Write-Host ''
      Show-Confetti
    }
    Write-Host ''
    Write-Host "  $(Icon 'Party')$($C.Bold)$($C.Green)linkwarden $installed is ready!$($C.Reset)"
    $first = Get-Command $Bin -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($first -and -not (Test-SameFolder (Split-Path -Parent $first.Path) $S.InstallDir)) {
      Write-Warn "Another $Bin comes first on your PATH: $(Format-Path $first.Path)"
      Write-Info "   Remove it or put $(Format-Path $S.InstallDir) earlier in PATH to use this one."
    }
    Write-Host ''
    Write-Host "  $($C.Bold)Next up:$($C.Reset)"
    Write-Host "    $($C.Dim)>$($C.Reset) $($C.Cyan)linkwarden init$($C.Reset)     $($C.Dim)# connect to your Linkwarden$($C.Reset)"
    Write-Host "    $($C.Dim)>$($C.Reset) $($C.Cyan)linkwarden --help$($C.Reset)   $($C.Dim)# look around$($C.Reset)"
    Write-Host ''
    Write-Host "  $($C.Dim)Happy bookmarking! $(Icon 'Bookmark')$($C.Reset)"
    Write-Host ''
  }

  # Removes PATH entries for the folder from the User PATH, and from the
  # Machine PATH when running as admin.
  function Remove-FromPath([string]$Folder) {
    foreach ($scope in @('User', 'Machine')) {
      $raw = Get-RawPath $scope
      $entries = @($raw -split ';' | Where-Object { $_ })
      $kept = @($entries | Where-Object { -not (Test-SameFolder $_ $Folder) })
      if ($kept.Count -eq $entries.Count) {
        continue
      }
      if ($scope -eq 'Machine' -and -not $S.Admin) {
        Write-Warn "$Folder is on the system PATH; remove it from an administrator terminal"
        continue
      }
      Set-RawPath $scope ($kept -join ';')
      Write-Ok "Removed $(Format-Path $Folder) from the $($scope.ToLowerInvariant()) PATH"
    }
    $env:Path = (@($env:Path -split ';' | Where-Object { $_ -and -not (Test-SameFolder $_ $Folder) })) -join ';'
  }

  # Removes the lines Add-ProfileLine wrote to the PowerShell profile.
  function Remove-ProfileLine {
    $profilePath = Get-ProfilePath
    if (-not (Test-Path -LiteralPath $profilePath)) {
      return
    }
    $lines = @(Get-Content -LiteralPath $profilePath)
    $kept = New-Object System.Collections.Generic.List[string]
    $removed = $false
    for ($i = 0; $i -lt $lines.Count; $i++) {
      if ($lines[$i] -eq $ProfileMarker -and $i + 1 -lt $lines.Count -and $lines[$i + 1] -like "*completions\$Bin.ps1*") {
        if ($kept.Count -gt 0 -and $kept[$kept.Count - 1] -eq '') {
          $kept.RemoveAt($kept.Count - 1)
        }
        $i++
        $removed = $true
        continue
      }
      $kept.Add($lines[$i])
    }
    if ($removed) {
      Set-Content -LiteralPath $profilePath -Value $kept -Encoding UTF8
      Write-Ok "Removed the completions line from $(Format-Path $profilePath)"
    }
  }

  # Removes linkwarden.exe from -Dir, or from the per-user and Program Files
  # folders, with its completions and the PATH and profile lines the
  # installer added.
  function Invoke-Uninstall {
    $folders = $(if ($Opt.Dir) { @(Resolve-InstallPath $Opt.Dir) } else { @($UserDir, $SystemDir) })
    $removed = $false
    Write-Host ''
    foreach ($folder in $folders) {
      $target = Join-Path $folder $Exe
      if (-not (Test-Path -LiteralPath $target)) {
        Write-Info "No $Exe at $(Format-Path $target)"
        continue
      }
      if (-not (Confirm-Choice "Remove $(Format-Path $target)?" $true)) {
        Write-Info 'Left it alone.'
        continue
      }
      try {
        Remove-Item -LiteralPath $target -Force
      } catch {
        Write-Warn "couldn't remove $(Format-Path $target): $(Get-ErrorText $_)"
        continue
      }
      Write-Ok "Removed $(Format-Path $target)"
      $completion = Get-CompletionPath $folder
      if (Test-Path -LiteralPath $completion) {
        Remove-Item -LiteralPath $completion -Force
        Write-Ok "Removed $(Format-Path $completion)"
      }
      foreach ($empty in @((Split-Path -Parent $completion), $folder)) {
        if ((Test-Path -LiteralPath $empty) -and -not (Get-ChildItem -LiteralPath $empty -Force)) {
          Remove-Item -LiteralPath $empty -Force
        }
      }
      Remove-FromPath $folder
      $removed = $true
    }
    Remove-ProfileLine
    $still = Get-Command $Bin -CommandType Application -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($still) {
      Write-Warn "$Bin is still on your PATH at $(Format-Path $still.Path)"
    }
    if ($removed) {
      Write-Host ''
      Write-Host "  $($C.Dim)Sorry to see you go. $(Icon 'Wave')$($C.Reset)"
      Write-Host ''
    }
  }

  $exitCode = 0
  try {
    if ($Opt.Help) {
      Show-Usage
      $S.Done = $true
      return 0
    }
    Initialize-Terminal

    if ($Opt.Uninstall) {
      Invoke-Uninstall
      $S.Done = $true
      return 0
    }

    Show-Intro
    $S.TmpDir = Join-Path ([IO.Path]::GetTempPath()) ("linkwarden-" + [Guid]::NewGuid().ToString('N'))
    New-Item -ItemType Directory -Path $S.TmpDir | Out-Null
    New-HttpClient

    Write-Step 1 'Sniffing out your system'
    Get-Platform
    Show-Platform

    Write-Step 2 'Finding the latest release'
    Resolve-Version
    if (-not (Test-Existing)) {
      $S.Done = $true
      return 0
    }

    Write-Step 3 'Picking a home'
    Select-InstallDir

    Write-Step 4 'Downloading'
    Get-Release

    Write-Step 5 "Checking it's the real deal"
    Test-Checksum

    Write-Step 6 'Moving in'
    Install-Binary
    Install-CompletionScript
    Add-ToPath

    Show-Finale
    $S.Done = $true
  } catch {
    $S.Failed = $true
    $exitCode = 1
    Clear-Line
    Write-Err ''
    Write-Err ("  {0}{1} {2}{3}" -f "$($C.Red)$($C.Bold)", $G.Fail, (Get-ErrorText $_), $C.Reset)
    Write-Err ''
  } finally {
    Show-Cursor
    if ($S.Client) {
      $S.Client.Dispose()
    }
    if ($S.TmpDir -and (Test-Path -LiteralPath $S.TmpDir)) {
      Remove-Item -LiteralPath $S.TmpDir -Recurse -Force -ErrorAction SilentlyContinue
    }
    # Neither finished nor failed: Ctrl-C stopped it.
    if (-not $S.Done -and -not $S.Failed) {
      Clear-Line
      Write-Err ''
      Write-Err "  $($C.Yellow)Install cancelled. Nothing was changed after this point. $(Icon 'Wave')$($C.Reset)"
      Write-Err ''
    }
  }
  return $exitCode
}

# Run from a saved file, the script exits with the status. Run through
# `irm | iex` or a script block, it can't: exit would close your terminal.
$lwFromFile = $false
try {
  $lwFromFile = "$($MyInvocation.MyCommand.ScriptContents)".Contains('function Install-LinkwardenCli')
} catch {
  $lwFromFile = $false
}
$lwExitCode = @(Install-LinkwardenCli -Opt @{
    Version       = "$Version".Trim().TrimStart('v')
    Dir           = $Dir
    Yes           = [bool]$Yes
    Force         = [bool]$Force
    NoCompletions = [bool]$NoCompletions
    NoModifyPath  = [bool]$NoModifyPath
    Plain         = [bool]$Plain
    Quiet         = [bool]$Quiet
    Uninstall     = [bool]$Uninstall
    Help          = [bool]$Help
    InSession     = -not $lwFromFile
  })[-1]
Remove-Item -Path Function:\Install-LinkwardenCli -ErrorAction SilentlyContinue
$global:LASTEXITCODE = $lwExitCode
if ($lwFromFile) {
  exit $lwExitCode
}
