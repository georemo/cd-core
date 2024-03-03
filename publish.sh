Version="v0.0.10"
# projDir="./cd-user/user.go"

# cd $projDir
go mod tidy
git submodule update --remote
git add -A
git commit -a -m "set version $Version"
git tag $Version
git push origin $Version

# cd-cli mod publish 

