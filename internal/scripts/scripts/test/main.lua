print("Hello from main.lua")
print("Output Directory:", GetOutputDir())
Templates:ShowTemplate("go", "default")
print("Has language Go: ", Templates:HasLanguage("go"))
print("Has language template go/default: ", Templates:HasTemplate("go", "default"))
print("Has language gozer: ", Templates:HasLanguage("gozer"))
print("Has language template gozer/busted: ", Templates:HasTemplate("gozer", "busted"))

local appConfigDir = GetAppConfigDir()
print("GetAppConfigDir: ", appConfigDir)
print("GetScriptsDir: ", GetScriptsDir())
print("GetTemplatesDir: ", GetTemplatesDir())

--success , errString = Downloader:GetFileWithProgress("20MB-TESTFILE.ORG.pdf", "https://files.testfile.org/PDF/20MB-TESTFILE.ORG.pdf")
--success, errString = Downloader:GetFileWithProgress("20MB-TESTFILE.ORG.pdf", "https://files.testfile.org/PDF/20MB-TESTFILE.ORG.pdferr")

-- Download a file from the internet that will WON"T work
--success, errString = Downloader:GetFile("20MB-TESTFILE.ORG.pdf", "https://files.testfile.org/PDF/20MB-TESTFILE.ORG.pdff")
-- Download a file from the internet that will work
success, errString = Downloader:GetFile("20MB-TESTFILE.ORG.pdf", "https://files.testfile.org/PDF/20MB-TESTFILE.ORG.pdf")

if not success then
    print("DownloadFile Error: ", errString)
end

local filepath = require("filepath")
local result = filepath.dir("/var/tmp/file.name")
print(result)

local shouldExit = Prompts:ConfirmationPrompt("Exit this script?")

if shouldExit then
    print("Exiting script")
else
    print("Not exiting script")
end

local inputtedValue = Prompts:InputPrompt("Input Test Value: ", "Default Value")
print("Inputted Value: ", inputtedValue)

local multilineInput = Prompts:MultiLineInputPrompt("Enter Multiline Value: ", "Default Value")
print("Multiline Input: ", multilineInput)

local choice = Prompts:ChoicePrompt("Choose a value", { "Option 1", "Option 2", "Option 3" }, "Option 2")
print("Choice: ", choice)

local choices = Prompts:MultiSelectPrompt("Favorite Animals?", { "Cat", "Dog", "Bird", "Whale" }, { "Cat", "Whale" })
for _, v in choices() do
    print(v)
end
