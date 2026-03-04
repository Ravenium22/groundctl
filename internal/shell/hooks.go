package shell

// BashHook returns a bash script that auto-checks on cd.
func BashHook() string {
	return `# groundctl shell hook for bash
# Add to ~/.bashrc: eval "$(ground hook bash)"

__groundctl_cd() {
    \builtin cd "$@" || return
    if [ -f .ground.yaml ]; then
        ground check --quiet 2>/dev/null
    fi
}
alias cd='__groundctl_cd'
`
}

// ZshHook returns a zsh script that auto-checks on chpwd.
func ZshHook() string {
	return `# groundctl shell hook for zsh
# Add to ~/.zshrc: eval "$(ground hook zsh)"

__groundctl_chpwd() {
    if [[ -f .ground.yaml ]]; then
        ground check --quiet 2>/dev/null
    fi
}
autoload -Uz add-zsh-hook
add-zsh-hook chpwd __groundctl_chpwd
`
}

// FishHook returns a fish script that auto-checks on directory change.
func FishHook() string {
	return `# groundctl shell hook for fish
# Add to ~/.config/fish/config.fish: ground hook fish | source

function __groundctl_cd --on-variable PWD
    if test -f .ground.yaml
        ground check --quiet 2>/dev/null
    end
end
`
}

// PowerShellHook returns a PowerShell script for auto-check on cd.
func PowerShellHook() string {
	return `# groundctl shell hook for PowerShell
# Add to $PROFILE: Invoke-Expression (ground hook powershell)

function Set-GroundctlLocation {
    param([string]$Path)
    Set-Location $Path
    if (Test-Path .ground.yaml) {
        ground check --quiet 2>$null
    }
}
Set-Alias -Name cd -Value Set-GroundctlLocation -Option AllScope -Force
`
}

// StarshipSegment returns a TOML config snippet for Starship prompt.
func StarshipSegment() string {
	return `# groundctl Starship prompt segment
# Add to ~/.config/starship.toml:

[custom.groundctl]
command = "ground check --quiet 2>/dev/null | tail -1 | tr -d ' '"
when = "test -f .ground.yaml"
format = "[$output]($style) "
style = "bold yellow"
description = "groundctl drift status"
`
}

// P10kSegment returns instructions for Powerlevel10k.
func P10kSegment() string {
	return `# groundctl Powerlevel10k prompt segment
# Add to ~/.p10k.zsh in prompt_custom_groundctl():

function prompt_groundctl() {
    local drift
    drift=$(ground check --json 2>/dev/null | grep -c '"status": "error"')
    if [[ "$drift" -gt 0 ]]; then
        p10k segment -f red -t "⚠ ${drift} drift"
    fi
}
# Then add 'groundctl' to POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS
`
}
