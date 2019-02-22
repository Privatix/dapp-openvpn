<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Disabled</key>
    <false/>
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
    </dict>
    <key>Label</key>
    <string>{{.Name}}</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.Script}}</string>
	<string>on</string>
	<string>{{.Server}}</string>
	<string>{{.Port}}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
