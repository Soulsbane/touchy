--Templates:CreateFileFromTemplate("d", "default", "app.d")
--Templates:List("all")
data, info = Templates:GetLanguageTemplateFor("go", "default")
print("Language Template for go/default: \n")
print(data)
print("Name: ", info.Name)
print("Description: ", info.Description)
print("DefaultOutputFileName: ", info.DefaultOutputFileName)
print("Embedded: ", info.Embedded)

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

local lsOutput, err = Command:RunWithOutput("ls")
print(lsOutput)
print("Command output with args: ")

local lsOutputWithArgs, err = Command:RunWithOutput("ls", "-al")
print(lsOutputWithArgs)
