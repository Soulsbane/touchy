--Templates:CreateFileFromTemplate("d", "default", "app.d")
--Templates:List("all")
data, info = Templates:GetLanguageTemplateFor("go", "default")
print("Language Template for go/default: \n")
print(data)
print(info.Name)
print(info.Description)
print(info.DefaultOutputFileName)
print(info.Embedded)

local createPathErr = IO:CreateDirAll("testpath/secondtestpath")

if createPathErr then
	print(createPathErr:Error())
else
	print("Directory path creation succeeded")
end

outputDirErr = IO:CreateDirInOutputDir("outputdir")

if outputDirErr then
	print(outputDirErr:Error())
else
	print("Output directory creation succeeded")
end
