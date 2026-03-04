# groundctl Demo — VHS tape file
# Record with: vhs demo.tape
# Requires: https://github.com/charmbracelet/vhs

Output demo.gif
Set FontSize 16
Set Width 900
Set Height 500
Set Theme "Catppuccin Mocha"
Set Padding 20

Type "# groundctl — terraform plan for your local dev machine"
Enter
Sleep 1s

Type "ground init"
Enter
Sleep 3s

Type "# Now let's check against the team standard"
Enter
Sleep 1s

Type "ground check"
Enter
Sleep 3s

Type "# Let's see what 'ground fix' would do"
Enter
Sleep 1s

Type "ground fix --dry-run"
Enter
Sleep 3s

Type "# Check the health of our setup"
Enter
Sleep 1s

Type "ground doctor"
Enter
Sleep 3s

Type "# That's it! 3 commands to manage your dev environment."
Enter
Sleep 2s
