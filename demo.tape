# groundctl Demo — VHS tape file
# Record with: vhs demo.tape
# Requires: https://github.com/charmbracelet/vhs

Output demo.gif
Set FontSize 16
Set Width 900
Set Height 550
Set Theme "Catppuccin Mocha"
Set Padding 20
Set Shell "powershell"

# Init
Type ".\ground.exe init"
Enter
Sleep 3s

# Check
Type ".\ground.exe check"
Enter
Sleep 3s

# Fix dry-run
Type ".\ground.exe fix --dry-run"
Enter
Sleep 3s

# Doctor
Type ".\ground.exe doctor"
Enter
Sleep 3s

Sleep 2s
