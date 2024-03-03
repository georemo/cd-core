Version="v0.0.2"
sourceFile="user.go"
projDir="~/cd-user/user.go"

cd $projDir
go mod tidy
git add $sourceFile
git commit -m "set version $Version"
git tag $Version
git push origin $Version

# cd-cli mod publish 

