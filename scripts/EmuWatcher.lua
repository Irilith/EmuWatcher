-- DISCLAIMER:
-- This script is created solely for educational purposes and does not condone or promote any illegal activities.
-- The use of this script is entirely at the user's discretion and risk.
-- The creator of this script (i) assumes no responsibility for any misuse or illegal actions resulting from its use.
-- Users are solely responsible for ensuring that their use of this script complies with all applicable laws and regulations.

local username = game:GetService("Players").LocalPlayer.Name
-- This script will write the current time to a file every 30 seconds
-- This is useful for checking if the user account is still running
-- If the file is not updated for a long time, the user may have left the game or the game may have crashed
while true do
	local unix = os.time()
	writefile("Emu.Watcher", username .. "|" .. unix)
	task.wait(30)
end
