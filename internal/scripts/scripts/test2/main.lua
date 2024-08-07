--Templates:CreateFileFromTemplate("d", "default", "app.d")
--Templates:List("all")
data, info = Templates:GetLanguageTemplateFor("go", "default")
print("Language Template for go/default: \n")
print(data)
print(info.Name)
print(info.Description)
print(info.DefaultOutputFileName)
print(info.Embedded)

local err = IO:CreateDirAll("test/path")

if err then
	print(err:Error())
else
	print("No error")
end
